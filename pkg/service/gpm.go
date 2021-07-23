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
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/hpcloud/tail"

	"github.com/lack-io/vine"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

type Gpm interface {
	Init() error
	ListService(context.Context) ([]*gpmv1.Service, int64, error)
	GetService(context.Context, string) (*gpmv1.Service, error)
	CreateService(context.Context, *gpmv1.Service) (*gpmv1.Service, error)
	StartService(context.Context, string) (*gpmv1.Service, error)
	StopService(context.Context, string) (*gpmv1.Service, error)
	RebootService(context.Context, string) (*gpmv1.Service, error)
	DeleteService(context.Context, string) (*gpmv1.Service, error)
	CatServiceLog(context.Context, string) ([]byte, error)
	WatchServiceLog(context.Context, string) (<-chan *gpmv1.ProcLog, error)

	Ls(context.Context, string) ([]*gpmv1.FileInfo, error)
	Pull(context.Context, string) (<-chan *gpmv1.PullResult, error)
	Push(context.Context, <-chan *gpmv1.PushIn) (<-chan *gpmv1.PushResult, error)
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

	return nil
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

func (g *gpm) CreateService(ctx context.Context, service *gpmv1.Service) (*gpmv1.Service, error) {
	if v, _ := g.GetService(ctx, service.Name); v != nil {
		return nil, verrs.Conflict(g.Name(), "service %s already exists", service.Name)
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

func (g *gpm) CatServiceLog(ctx context.Context, name string) ([]byte, error) {
	f := filepath.Join(g.Cfg.Root, "logs", name)
	stat, _ := os.Stat(f)
	if stat == nil {
		return nil, verrs.NotFound(g.Name(), "service '%s' log not exists", name)
	}

	var out []byte
	ech := make(chan error, 1)
	done := make(chan struct{}, 1)

	go func() {
		var err error
		out, err = ioutil.ReadFile(filepath.Join(f, name+".log"))
		if err != nil {
			ech <- err
			return
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return out, nil
	case e := <-ech:
		return nil, e
	case <-done:
		return out, nil
	}
}

func (g *gpm) WatchServiceLog(ctx context.Context, name string) (<-chan *gpmv1.ProcLog, error) {
	f := filepath.Join(g.Cfg.Root, "logs", name)
	stat, _ := os.Stat(f)
	if stat == nil {
		return nil, verrs.NotFound(g.Name(), "service '%s' log not exists", name)
	}

	out := make(chan *gpmv1.ProcLog, 10)
	ech := make(chan error, 1)

	go func() {
		cfg := tail.Config{
			Poll:   true,
			Follow: true,
		}

		t, err := tail.TailFile(filepath.Join(f, name+".log"), cfg)
		if err != nil {
			ech <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case line := <-t.Lines:
				l := &gpmv1.ProcLog{
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
