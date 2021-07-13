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
	"time"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/lack-io/vine"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

func init() {
	inject.ProvidePanic(new(gpm))
}

type Gpm interface {
	Init() error
	ListService(context.Context, *gpmv1.PageMeta) ([]*gpmv1.Service, int64, error)
	GetService(context.Context, int64) (*gpmv1.Service, error)
	GetServiceByName(context.Context, string) (*gpmv1.Service, error)
	CreateService(context.Context, *gpmv1.Service) (*gpmv1.Service, error)
	StartService(context.Context, int64) (*gpmv1.Service, error)
	StopService(context.Context, int64) (*gpmv1.Service, error)
	RebootService(context.Context, int64) (*gpmv1.Service, error)
	DeleteService(context.Context, int64) (*gpmv1.Service, error)
}

var _ Gpm = (*gpm)(nil)

type gpm struct {
	vine.Service `inject:""`
}

func (g *gpm) Init() error {
	return nil
}

func (g *gpm) ListService(ctx context.Context, meta *gpmv1.PageMeta) ([]*gpmv1.Service, int64, error) {
	outs, total, err := dao.ServiceSBuilder().FindPage(ctx, int(meta.Page), int(meta.Size))
	if err != nil {
		return nil, 0, err
	}

	//for _, item := range outs {
	//
	//}

	return outs, total, nil
}

func (g *gpm) GetService(ctx context.Context, id int64) (*gpmv1.Service, error) {
	s, err := dao.ServiceSBuilder().SetId(id).FindOne(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (g *gpm) GetServiceByName(ctx context.Context, name string) (*gpmv1.Service, error) {
	s, err := dao.ServiceSBuilder().SetName(name).FindOne(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (g *gpm) CreateService(ctx context.Context, service *gpmv1.Service) (*gpmv1.Service, error) {
	if v, _ := g.GetServiceByName(ctx, service.Name); v != nil {
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

	service, err = dao.FromService(service).Create(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (g *gpm) StartService(ctx context.Context, id int64) (*gpmv1.Service, error) {
	panic("implement me")
}

func (g *gpm) StopService(ctx context.Context, id int64) (*gpmv1.Service, error) {
	panic("implement me")
}

func (g *gpm) RebootService(ctx context.Context, id int64) (*gpmv1.Service, error) {
	panic("implement me")
}

func (g *gpm) DeleteService(ctx context.Context, id int64) (*gpmv1.Service, error) {
	panic("implement me")
}
