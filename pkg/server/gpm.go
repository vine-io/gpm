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
	gruntime "runtime"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	"github.com/gpm2/gpm/pkg/service"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/cli"
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
	rsp.Services, rsp.Total, err = s.H.ListService(ctx)
	return
}

func (s *server) GetService(ctx context.Context, req *pb.GetServiceReq, rsp *pb.GetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.GetService(ctx, req.Name)
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
		Dir:         req.Dir,
		Env:         req.Env,
		SysProcAttr: req.SysProcAttr,
		Log:         req.Log,
		Version:     req.Version,
		AutoRestart: req.AutoRestart,
	}
	rsp.Service, err = s.H.CreateService(ctx, ss)
	return
}

func (s *server) StartService(ctx context.Context, req *pb.StartServiceReq, rsp *pb.StartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StartService(ctx, req.Name)
	return
}

func (s *server) StopService(ctx context.Context, req *pb.StopServiceReq, rsp *pb.StopServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StopService(ctx, req.Name)
	return
}

func (s *server) RebootService(ctx context.Context, req *pb.RebootServiceReq, rsp *pb.RebootServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.RebootService(ctx, req.Name)
	return
}

func (s *server) DeleteService(ctx context.Context, req *pb.DeleteServiceReq, rsp *pb.DeleteServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadGateway(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.DeleteService(ctx, req.Name)
	return
}

func (s *server) CatServiceLog(ctx context.Context, req *pb.CatServiceLogReq, rsp *pb.CatServiceLogRsp) error {
	panic("implement me")
}

func (s *server) WatchServiceLog(ctx context.Context, req *pb.WatchServiceLogReq, stream pb.GpmService_WatchServiceLogStream) error {
	panic("implement me")
}

func (s *server) InstallService(ctx context.Context, stream pb.GpmService_InstallServiceStream) error {
	panic("implement me")
}

func (s *server) ListServiceVersions(ctx context.Context, req *pb.ListServiceVersionsReq, rsp *pb.ListServiceVersionsRsp) error {
	panic("implement me")
}

func (s *server) UpgradeService(ctx context.Context, stream pb.GpmService_UpgradeServiceStream) error {
	panic("implement me")
}

func (s *server) RollBackService(ctx context.Context, req *pb.RollbackServiceReq, stream pb.GpmService_RollBackServiceStream) error {
	panic("implement me")
}

func (s *server) Ls(ctx context.Context, req *pb.LsReq, rsp *pb.LsRsp) error {
	panic("implement me")
}

func (s *server) Pull(ctx context.Context, req *pb.PullReq, stream pb.GpmService_PullStream) error {
	panic("implement me")
}

func (s *server) Push(ctx context.Context, stream pb.GpmService_PushStream) error {
	panic("implement me")
}

func (s *server) Exec(ctx context.Context, req *pb.ExecReq, rsp *pb.ExecRsp) error {
	panic("implement me")
}

func (s *server) Terminal(ctx context.Context, stream pb.GpmService_TerminalStream) error {
	panic("implement me")
}

func (s *server) Init() error {
	var err error

	opts := []vine.Option{
		vine.Name(runtime.GpmName),
		vine.Id(runtime.GpmId),
		vine.Version(runtime.GetVersion()),
		vine.Metadata(map[string]string{
			"namespace": runtime.Namespace,
		}),
		vine.Flags(&cli.StringFlag{
			Name:    "root",
			Usage:   "gpmd root directory",
			EnvVars: []string{"GPMD_ROOT"},
		}),
		vine.Action(func(c *cli.Context) error {

			cfg := &config.Config{}
			cfg.Root = c.String("root")
			if cfg.Root == "" {
				if gruntime.GOOS == "windows" {
					cfg.Root = "C:\\opt\\lack\\gpmd"
				} else {
					cfg.Root = "/opt/lack/gpmd"
				}
			}

			return inject.Provide(cfg)
		}),
	}

	s.Service.Init(opts...)

	db := new(dao.DB)
	if err = inject.Provide(s.Service, s.Client(), s, db); err != nil {
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
