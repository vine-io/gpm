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

package biz

import (
	"context"
	"io"
	"os"
	"path/filepath"
	gruntime "runtime"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/shirou/gopsutil/mem"
	proc "github.com/shirou/gopsutil/process"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/infra/repo"
	"github.com/vine-io/gpm/pkg/runtime"
	"github.com/vine-io/gpm/pkg/runtime/inject"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/lib/config"
	verrs "github.com/vine-io/vine/lib/errors"
	log "github.com/vine-io/vine/lib/logger"
)

func init() {
	inject.ProvidePanic(new(manager))
}

type manager struct {
	vine.Service `inject:""`

	DB *repo.DB `inject:""`

	up time.Time
	sync.RWMutex
	ps map[string]*Process
}

func (g *manager) Init() error {
	var err error

	if err = os.MkdirAll(filepath.Join(config.Get("root").String(""), "services"), 0o777); err != nil {
		return err
	}
	ctx := context.Background()
	list, err := g.DB.FindAllServices(ctx)
	if err != nil {
		return err
	}

	g.ps = map[string]*Process{}
	for _, item := range list {
		p := NewProcess(item)
		if item.Status == gpmv1.StatusRunning {
			_, _ = g.startService(ctx, p)
		}
		g.ps[item.Name] = p
	}

	g.up = time.Now()
	return nil
}

func (g *manager) Info(ctx context.Context) (*gpmv1.GpmInfo, error) {
	info := &gpmv1.GpmInfo{
		Version: runtime.GetVersion(),
		Goos:    gruntime.GOOS,
		Arch:    gruntime.GOARCH,
		Gov:     gruntime.Version(),
		Pid:     int32(os.Getpid()),
	}

	stat := &gpmv1.Stat{}
	pr, _ := proc.NewProcess(info.Pid)
	if pr != nil {
		percent, _ := pr.MemoryPercent()
		stat.MemPercent = percent
		m, _ := mem.VirtualMemory()
		if m != nil {
			stat.Memory = uint64(float64(percent) / 100 * float64(m.Total))
		}
		stat.CpuPercent, _ = pr.CPUPercent()
	}
	info.Stat = stat
	info.UpTime = time.Now().Unix() - g.up.Unix()

	return info, nil
}

func (g *manager) List(ctx context.Context) ([]*gpmv1.Service, int64, error) {
	outs, err := g.DB.FindAllServices(ctx)
	if err != nil {
		return nil, 0, err
	}

	for i := 0; i < len(outs); i++ {
		statProcess(outs[i])
	}

	return outs, int64(len(outs)), nil
}

func (g *manager) Get(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	statProcess(s)

	return s, nil
}

func (g *manager) getService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.DB.FindService(ctx, name)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (g *manager) Create(ctx context.Context, spec *gpmv1.ServiceSpec) (*gpmv1.Service, error) {
	if v, _ := g.getService(ctx, spec.Name); v != nil {
		return nil, verrs.Conflict(g.Name(), "service %s already exists", spec.Name)
	}

	service := &gpmv1.Service{
		Name:        spec.Name,
		Bin:         spec.Bin,
		Args:        spec.Args,
		Dir:         spec.Dir,
		Env:         spec.Env,
		SysProcAttr: spec.SysProcAttr,
		Log:         spec.Log,
		Version:     spec.Version,
		AutoRestart: spec.AutoRestart,
		InstallFlag: spec.InstallFlag,
	}

	err := fillService(service)
	if err != nil {
		return nil, err
	}

	service.Status = gpmv1.StatusInit
	service.CreationTimestamp = time.Now().Unix()
	if service.Version == "" {
		service.Version = "v0.0.1"
	}

	service, err = g.DB.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}

	g.Lock()
	g.ps[service.Name] = NewProcess(service)
	g.Unlock()

	return service, nil
}

func (g *manager) Edit(ctx context.Context, name string, spec *gpmv1.EditServiceSpec) (*gpmv1.Service, error) {
	service, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p, ok := g.ps[name]
	g.RUnlock()

	var isRunning bool
	if ok && p.Status == gpmv1.StatusRunning {
		isRunning = true
		g.stopService(ctx, p)
	}

	if spec.Bin != "" {
		service.Bin = spec.Bin
	}
	if len(spec.Env) > 0 {
		service.Env = spec.Env
	}
	if spec.Dir != "" {
		service.Dir = spec.Dir
	}
	if spec.Log != nil {
		service.Log = &gpmv1.ProcLog{}
		if spec.Log.Expire > 0 {
			service.Log.Expire = spec.Log.Expire
		}
		if spec.Log.MaxSize > 0 {
			service.Log.MaxSize = spec.Log.MaxSize
		}
	}
	if spec.SysProcAttr != nil {
		service.SysProcAttr = spec.SysProcAttr
	}
	if len(spec.Args) > 0 {
		service.Args = spec.Args
	}
	if spec.AutoRestart != 0 {
		service.AutoRestart = spec.AutoRestart
	}

	err = fillService(service)
	if err != nil {
		return nil, err
	}

	service, err = g.DB.UpdateService(ctx, service)
	if err != nil {
		return nil, err
	}

	p = NewProcess(service)
	if isRunning {
		g.startService(ctx, p)
	}

	g.Lock()
	g.ps[service.Name] = p
	g.Unlock()

	return service, nil
}

func (g *manager) Start(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()

	s, err = g.startService(ctx, p)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (g *manager) startService(ctx context.Context, p *Process) (*gpmv1.Service, error) {
	pid, err := p.Start()

	s := p.Service
	now := time.Now()
	s.UpdateTimestamp = now.Unix()
	if err != nil {
		s.Status = gpmv1.StatusFailed
		s.Msg = err.Error()
	} else {
		s.Pid = int64(pid)
		s.Status = gpmv1.StatusRunning
	}

	var e error
	s, e = g.DB.UpdateService(ctx, s)
	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, e
	}

	g.Lock()
	g.ps[s.Name] = p
	g.Unlock()

	return s, nil
}

func (g *manager) Stop(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()

	return g.stopService(ctx, p)
}

func (g *manager) stopService(ctx context.Context, p *Process) (*gpmv1.Service, error) {

	err := p.Stop()
	if err != nil {
		return nil, err
	}

	s := p.Service
	s.Status = gpmv1.StatusStopped
	now := time.Now()
	s.UpdateTimestamp = now.Unix()
	s, err = g.DB.UpdateService(ctx, s)
	if err != nil {
		return nil, err
	}

	g.Lock()
	g.ps[s.Name] = p
	g.Unlock()

	return s, nil
}

func (g *manager) Restart(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()

	g.stopService(ctx, p)

	return g.startService(ctx, p)
}

func (g *manager) Delete(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.getService(ctx, name)
	if err != nil {
		return nil, err
	}

	if s.Status == gpmv1.StatusRunning {

		g.RLock()
		p := g.ps[s.Name]
		g.RUnlock()

		if _, err = g.stopService(ctx, p); err != nil {
			return nil, err
		}
	}

	g.Lock()
	delete(g.ps, s.Name)
	g.Unlock()

	err = g.DB.DeleteService(ctx, s.Name)
	if err != nil {
		return nil, err
	}

	if s.InstallFlag == 1 {
		log.Infof("remove %s directory", s.Name)
		link, _ := os.Readlink(s.Dir)
		if link != "" {
			_ = os.RemoveAll(link)
		}
		_ = os.RemoveAll(s.Dir)
	}
	return s, nil
}

func (g *manager) TailLog(ctx context.Context, name string, number int64, follow bool, sender IOWriter) error {
	f := filepath.Join(config.Get("root").String(""), "logs", name, name+".log")
	stat, _ := os.Stat(f)
	if stat == nil {
		return verrs.NotFound(g.Name(), "service '%s' log not exists", name)
	}

	cfg := tail.Config{
		Poll: true,
	}

	if number > 0 {
		total := stat.Size()
		if number > total {
			cfg.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekStart}
		} else {
			cfg.Location = &tail.SeekInfo{Offset: -1 * number, Whence: io.SeekEnd}
		}
	}

	if follow {
		cfg.ReOpen = true
		cfg.Follow = true
	}

	t, err := tail.TailFile(f, cfg)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return nil
		case line, ok := <-t.Lines:
			if !ok {
				continue
			}
			l := &gpmv1.ServiceLog{
				Text:      line.Text,
				Timestamp: line.Time.Unix(),
			}
			if line.Err != nil {
				l.Error = line.Err.Error()
			}
			_ = sender.Send(l)
		}
	}
}

func (g *manager) String() string {
	return "manager"
}
