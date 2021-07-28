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
	"io"
	"os"
	"path/filepath"
	gruntime "runtime"
	"sync"
	"time"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/hpcloud/tail"
	"github.com/lack-io/vine"
	verrs "github.com/lack-io/vine/proto/apis/errors"
	"github.com/shirou/gopsutil/mem"
	proc "github.com/shirou/gopsutil/process"
)

type Gpm interface {
	Init() error
	Info(context.Context) (*gpmv1.GpmInfo, error)
	ListService(context.Context) ([]*gpmv1.Service, int64, error)
	GetService(context.Context, string) (*gpmv1.Service, error)
	CreateService(context.Context, *gpmv1.ServiceSpec) (*gpmv1.Service, error)
	StartService(context.Context, string) (*gpmv1.Service, error)
	StopService(context.Context, string) (*gpmv1.Service, error)
	RebootService(context.Context, string) (*gpmv1.Service, error)
	DeleteService(context.Context, string) (*gpmv1.Service, error)
	WatchServiceLog(context.Context, string, int64, bool) (<-chan *gpmv1.ServiceLog, error)

	InstallService(context.Context, *gpmv1.ServiceSpec, <-chan *gpmv1.Package) (<-chan *gpmv1.InstallServiceResult, error)
	ListServiceVersions(context.Context, string) ([]*gpmv1.ServiceVersion, error)
	UpgradeService(context.Context, string, string, <-chan *gpmv1.Package) (<-chan *gpmv1.UpgradeServiceResult, error)
	RollbackService(context.Context, string, string) error

	Ls(context.Context, string) ([]*gpmv1.FileInfo, error)
	Pull(context.Context, string, bool) (<-chan *gpmv1.PullResult, error)
	Push(context.Context, string, string, <-chan *gpmv1.PushIn) (<-chan *gpmv1.PushResult, error)
	Exec(context.Context, *gpmv1.ExecIn) (<-chan *gpmv1.ExecResult, error)
	Terminal(context.Context, <-chan *gpmv1.TerminalIn) (<-chan *gpmv1.TerminalResult, error)
}

func init() {
	inject.ProvidePanic(new(gpm))
}

var _ Gpm = (*gpm)(nil)

type gpm struct {
	vine.Service `inject:""`

	Cfg *config.Config `inject:""`
	DB  *dao.DB        `inject:""`

	up time.Time
	sync.RWMutex
	ps map[string]*Process
}

func (g *gpm) Init() error {
	var err error

	if err = os.MkdirAll(filepath.Join(g.Cfg.Root, "services"), os.ModePerm); err != nil {
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

func (g *gpm) Info(ctx context.Context) (*gpmv1.GpmInfo, error) {
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

func (g *gpm) ListService(ctx context.Context) ([]*gpmv1.Service, int64, error) {
	outs, err := g.DB.FindAllServices(ctx)
	if err != nil {
		return nil, 0, err
	}

	for i := 0; i < len(outs); i++ {
		statProcess(outs[i])
	}

	return outs, int64(len(outs)), nil
}

func (g *gpm) GetService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.DB.FindService(ctx, name)
	if err != nil {
		return nil, err
	}

	statProcess(s)

	return s, nil
}

func (g *gpm) CreateService(ctx context.Context, spec *gpmv1.ServiceSpec) (*gpmv1.Service, error) {
	if v, _ := g.GetService(ctx, spec.Name); v != nil {
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

func (g *gpm) StartService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.GetService(ctx, name)
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

func (g *gpm) startService(ctx context.Context, p *Process) (*gpmv1.Service, error) {
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
		s.StartTimestamp = now.Unix()
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

func (g *gpm) StopService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.GetService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()

	return g.stopService(ctx, p)
}

func (g *gpm) stopService(ctx context.Context, p *Process) (*gpmv1.Service, error) {

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

func (g *gpm) RebootService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.GetService(ctx, name)
	if err != nil {
		return nil, err
	}

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()

	g.stopService(ctx, p)

	return g.startService(ctx, p)
}

func (g *gpm) DeleteService(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := g.GetService(ctx, name)
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
	return s, err
}

func (g *gpm) WatchServiceLog(ctx context.Context, name string, number int64, follow bool) (<-chan *gpmv1.ServiceLog, error) {
	f := filepath.Join(g.Cfg.Root, "logs", name, name+".log")
	stat, _ := os.Stat(f)
	if stat == nil {
		return nil, verrs.NotFound(g.Name(), "service '%s' log not exists", name)
	}

	out := make(chan *gpmv1.ServiceLog, 10)
	ech := make(chan error, 1)

	go func() {
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
			ech <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
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
				out <- l
			}
		}
	}()

	return out, nil
}
