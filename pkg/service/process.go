// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/lack-io/pkg/unit"
	log "github.com/lack-io/vine/lib/logger"
	"github.com/shirou/gopsutil/mem"
	proc "github.com/shirou/gopsutil/process"
)

const timeFormat = "20060102150405"

var (
	ErrProcessNotFound = errors.New("process not found")
)

type Process struct {
	*gpmv1.Service

	pr *proc.Process

	cfg *config.Config
	db  *dao.DB

	lw io.WriteCloser

	done chan struct{}
}

func NewProcess(in *gpmv1.Service) *Process {
	process := &Process{
		Service: in,
		cfg:     &config.Config{},
		db:      &dao.DB{},
		done:    make(chan struct{}, 1),
	}
	if process.Pid != 0 {
		p, err := proc.NewProcess(int32(process.Pid))
		if err == nil {
			process.pr = p
		} else {
			process.Pid = 0
		}
	}

	_ = inject.Resolve(process.cfg)
	_ = inject.Resolve(process.db)

	return process
}

func (p *Process) Start() (int32, error) {
	var pid int32
	var err error
	if p.pr == nil || p.pr.Pid == 0 {
		pid, err = p.run()
		if err != nil {
			return 0, err
		}
	}
	pid = int32(p.Pid)

	p.done = make(chan struct{}, 1)
	p.done <- struct{}{}

	if p.AutoRestart > 0 {
		go p.watching()
	}
	if p.Log != nil {
		go p.rotating()
	}

	return pid, nil
}

func (p *Process) run() (int32, error) {
	cmd := exec.Command(p.Bin, p.Args...)

	if p.Env != nil {
		env := make([]string, 0)
		for k, v := range p.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	if p.Dir != "" {
		cmd.Dir = p.Dir
	}

	if p.SysProcAttr != nil {
		injectSysProcAttr(cmd, p.SysProcAttr)
	}

	root := filepath.Join(p.cfg.Root, "logs", p.Name)
	_ = os.MkdirAll(root, os.ModePerm)

	flog := filepath.Join(root, p.Name+".log")
	_ = os.Rename(flog, filepath.Join(root, fmt.Sprintf("%s.log-%s", p.Name, time.Now().Format(timeFormat))))

	var err error
	p.lw, err = os.OpenFile(flog, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return 0, err
	}

	cmd.Stdout = p.lw
	cmd.Stderr = p.lw

	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	var pr *os.Process
	for {
		pr = cmd.Process
		if pr != nil {
			break
		}
		time.Sleep(time.Millisecond * 300)
	}

	p.pr, _ = proc.NewProcess(int32(pr.Pid))
	p.Pid = int64(pr.Pid)
	return p.pr.Pid, nil
}

func (p *Process) watching() {
	log.Infof("start service %s(%d) watching", p.Name, p.Pid)
	timer := time.NewTicker(time.Second * 2)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-p.done:
			if !ok {
				log.Infof("stop service %s(%d) watching", p.Name, p.Pid)
				return
			}
		case <-timer.C:
			pr, err := proc.NewProcess(int32(p.Pid))
			if err != nil {
				pid, _ := p.run()
				log.Infof("reboot service(dead) %s at pid: %d", p.Name, pid)
				p.db.UpdateService(context.TODO(), p.Service)
			} else {
				status, _ := pr.Status()
				if status == "Z" {
					log.Infof("watching service(pid=%d) %s status: %s", p.Pid, p.Name, status)
					_ = p.kill()
					pid, _ := p.run()
					p.db.UpdateService(context.TODO(), p.Service)
					log.Infof("reboot service(Z) %s at pid: %d", p.Name, pid)
				}
			}
		}
	}
}

func (p *Process) rotating() {
	timer := time.NewTicker(time.Second * 2)
	log.Infof("start service %s(%d) rotating", p.Name, p.Pid)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-p.done:
			if !ok {
				log.Infof("stop service %s(%d) rotating", p.Name, p.Pid)
				return
			}
		case <-timer.C:
			now := time.Now()
			param := p.Log
			// 日志目录
			root := filepath.Join(p.cfg.Root, "logs", p.Name)
			// 当前日志文件
			plog := filepath.Join(root, p.Name+".log")

			stat, _ := os.Stat(plog)
			// 日志大小超过额定值，进行日志切分
			if stat != nil && stat.Size() > param.MaxSize {
				log.Infof("log %s greater than %s, rotates it", plog, unit.ConvAuto(param.MaxSize, 2))
				err := rotate(plog, stat.Size(), param.MaxSize)
				if err != nil {
					log.Errorf("%v", err)
				}
			}

			// 遍历服务所有日志文件
			filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && root != path {
					return filepath.SkipDir
				}

				if !d.Type().IsRegular() {
					return nil
				}

				name := d.Name()
				parts := strings.Split(name, "-")
				if len(parts) > 1 {
					logT, _ := time.Parse(timeFormat, parts[1])
					// 删除过期的日志文件
					if now.Sub(logT).Hours() > float64(param.Expire*24) {
						_ = os.Remove(path)
						log.Infof("remove expired log: %s", path)
					}
				}

				return nil
			})
		}
	}
}

func (p *Process) Kill() error {
	if p.pr == nil {
		return ErrProcessNotFound
	}

	close(p.done)

	return p.kill()
}

func (p *Process) kill() error {
	pr, err := os.FindProcess(int(p.pr.Pid))
	if err != nil {
		return err
	}
	err = pr.Kill()
	if err != nil {
		return err
	}
	_, err = pr.Wait()
	if err != nil {
		return err
	}
	if err = pr.Release(); err != nil {
		return err
	}
	p.Pid = 0
	p.pr = nil

	if p.lw != nil {
		return p.lw.Close()
	}
	return nil
}

func (p *Process) Stop() error {
	if p.pr == nil {
		return ErrProcessNotFound
	}

	close(p.done)

	return p.stop()
}

func (p *Process) stop() error {
	pr, err := os.FindProcess(int(p.pr.Pid))
	if err != nil {
		return err
	}
	err = pr.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	_, _ = pr.Wait()
	_ = pr.Release()
	p.Pid = 0
	p.pr = nil

	if p.lw != nil {
		return p.lw.Close()
	}
	return nil
}

func statProcess(s *gpmv1.Service) {
	var pr *proc.Process
	if s.Pid > 0 {
		pr, _ = proc.NewProcess(int32(s.Pid))
	}
	stat := &gpmv1.Stat{}
	if pr != nil {
		percent, _ := pr.MemoryPercent()
		stat.MemPercent = percent
		m, _ := mem.VirtualMemory()
		if m != nil {
			stat.Memory = uint64(float64(percent) / 100 * float64(m.Total))
		}
		stat.CpuPercent, _ = pr.CPUPercent()
	}
	s.Stat = stat
}

func rotate(rl string, total, size int64) error {
	pf, e := os.OpenFile(rl, os.O_RDWR|os.O_SYNC, os.ModePerm)
	if e != nil {
		return fmt.Errorf("open log file %s: %v", rl, e)
	}
	defer pf.Close()

	//_, e = pf.Seek(size, 0)
	buf := make([]byte, size)
	n, e := pf.Read(buf)
	if e != nil {
		return fmt.Errorf("open log file %s: %v", rl, e)
	}

	flog := rl + "-" + time.Now().Format(timeFormat)
	e = ioutil.WriteFile(flog, buf[:n], os.ModePerm)
	if e != nil {
		return fmt.Errorf("write log %s: %v", flog, e)
	}

	buf1 := make([]byte, total-size)
	n, e = pf.Read(buf1)
	if e != nil {
		return fmt.Errorf("open log file %s: %v", rl, e)
	}

	e = os.WriteFile(rl, buf1[:n], os.ModePerm)
	if e != nil {
		return fmt.Errorf("rewrite log file %s: %v", rl, e)
	}

	return nil
}
