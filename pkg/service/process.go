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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	log "github.com/lack-io/vine/lib/logger"
	"github.com/shirou/gopsutil/mem"
	proc "github.com/shirou/gopsutil/process"
)

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

	now := time.Now()
	root := filepath.Join(p.cfg.Root, "logs", p.Name)
	_ = os.MkdirAll(root, os.ModePerm)
	flog := filepath.Join(root, fmt.Sprintf("%s.log-%s", p.Name, now.Format("20060102150405")))

	err := ioutil.WriteFile(flog, []byte(""), os.ModePerm)
	if err != nil {
		return 0, nil
	}

	link := filepath.Join(root, p.Name+".log")
	_ = os.Remove(link)
	err = os.Symlink(flog, link)
	if err != nil {
		return 0, fmt.Errorf("%w: create soft link", err)
	}

	p.lw, err = os.OpenFile(link, os.O_RDWR|os.O_TRUNC, os.ModePerm)
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
	timer := time.NewTicker(time.Second * 2)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-p.done:
			if !ok {
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
