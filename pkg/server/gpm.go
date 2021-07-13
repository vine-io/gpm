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

package server

import (
	"context"

	"github.com/gpm2/gpm/pkg/runtime/inject"
	"github.com/gpm2/gpm/pkg/service"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/vine"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

type server struct {
	vine.Service

	H service.Gpm `inject:""`
}

func (s *server) Healthz(ctx context.Context, req *pb.Empty, rsp *pb.Empty) error {
	return nil
}

func (s *server) ListService(ctx context.Context, req *pb.ListServiceReq, rsp *pb.ListServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Servicees, rsp.Total, err = s.H.ListService(ctx, &req.PageMeta)
	return
}

func (s *server) GetService(ctx context.Context, req *pb.GetServiceReq, rsp *pb.GetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.GetService(ctx, req.Id)
	return
}

func (s *server) GetServiceByName(ctx context.Context, req *pb.GetServiceByNameReq, rsp *pb.GetServiceByNameRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.GetServiceByName(ctx, req.Name)
	return
}

func (s *server) CreateService(ctx context.Context, req *pb.CreateServiceReq, rsp *pb.CreateServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	ss := &gpmv1.Service{
		Name:        req.Name,
		Bin:         req.Bin,
		Args:        req.Args,
		Chroot:      req.Dir,
		User:        req.User,
		Group:       req.Group,
		Version:     req.Group,
		AutoRestart: req.AutoRestart,
	}
	rsp.Service, err = s.H.CreateService(ctx, ss)
	return
}

func (s *server) StartService(ctx context.Context, req *pb.StartServiceReq, rsp *pb.StartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StartService(ctx, req.Id)
	return
}

func (s *server) StopService(ctx context.Context, req *pb.StopServiceReq, rsp *pb.StopServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StopService(ctx, req.Id)
	return
}

func (s *server) RebootService(ctx context.Context, req *pb.RebootServiceReq, rsp *pb.RebootServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.RebootService(ctx, req.Id)
	return
}

func (s *server) DeleteService(ctx context.Context, req *pb.DeleteServiceReq, rsp *pb.DeleteServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.DeleteService(ctx, req.Id)
	return
}

func (s *server) Init(opts ...vine.Option) error {
	var err error
	s.Service.Init(opts...)

	if err = inject.Provide(s.Service, s.Client(), s); err != nil {
		return err
	}

	// TODO: inject more objects

	if err = inject.Populate(); err != nil {
		return err
	}

	if err = s.H.Init(); err != nil {
		return err
	}

	if err = pb.RegisterGpmServiceHandler(s.Service.Server(), s); err != nil {
		return err
	}

	return err
}

func New(opts ...vine.Option) *server {
	srv := vine.NewService(opts...)
	return &server{
		Service: srv,
	}
}
