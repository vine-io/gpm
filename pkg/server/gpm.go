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

	pb "github.com/vine-io/gpm/api/service/gpm/v1"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/service"
	vserver "github.com/vine-io/vine/core/server"
	verrs "github.com/vine-io/vine/lib/errors"
)

type GpmServer struct {
	server vserver.Server

	manager service.GenerateManager
	ftp     service.GenerateFTP
}

func RegistryGpmRpcServer(ctx context.Context, s vserver.Server, manager service.GenerateManager, ftp service.GenerateFTP) error {

	gs := &GpmServer{
		server:  s,
		manager: manager,
		ftp:     ftp,
	}

	if err := pb.RegisterGpmServiceHandler(s, gs); err != nil {
		return err
	}

	return nil
}

func (s *GpmServer) Name() string {
	return s.server.Options().Name
}

func (s *GpmServer) Healthz(ctx context.Context, _ *pb.Empty, rsp *pb.Empty) error {
	return nil
}

func (s *GpmServer) UpdateSelf(ctx context.Context, stream pb.GpmService_UpdateSelfStream) error {
	return s.manager.Update(ctx, &simpleUpdateSelfStream{stream: stream})
}

func (s *GpmServer) Info(ctx context.Context, _ *pb.InfoReq, rsp *pb.InfoRsp) (err error) {
	rsp.Gpm, err = s.manager.Info(ctx)
	return
}

func (s *GpmServer) ListService(ctx context.Context, req *pb.ListServiceReq, rsp *pb.ListServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Services, rsp.Total, err = s.manager.List(ctx)
	return
}

func (s *GpmServer) GetService(ctx context.Context, req *pb.GetServiceReq, rsp *pb.GetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Get(ctx, req.Name)
	return
}

func (s *GpmServer) CreateService(ctx context.Context, req *pb.CreateServiceReq, rsp *pb.CreateServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Create(ctx, req.Spec)
	return
}

func (s *GpmServer) EditService(ctx context.Context, req *pb.EditServiceReq, rsp *pb.EditServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Edit(ctx, req.Name, req.Spec)
	return
}

func (s *GpmServer) StartService(ctx context.Context, req *pb.StartServiceReq, rsp *pb.StartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Start(ctx, req.Name)
	return
}

func (s *GpmServer) StopService(ctx context.Context, req *pb.StopServiceReq, rsp *pb.StopServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Stop(ctx, req.Name)
	return
}

func (s *GpmServer) RestartService(ctx context.Context, req *pb.RestartServiceReq, rsp *pb.RestartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Restart(ctx, req.Name)
	return
}

func (s *GpmServer) DeleteService(ctx context.Context, req *pb.DeleteServiceReq, rsp *pb.DeleteServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.manager.Delete(ctx, req.Name)
	return
}

func (s *GpmServer) WatchServiceLog(ctx context.Context, req *pb.WatchServiceLogReq, stream pb.GpmService_WatchServiceLogStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	return s.manager.TailLog(ctx, req.Name, req.Number, req.Follow, &simpleWatchLogSender{stream: stream})
}

func (s *GpmServer) InstallService(ctx context.Context, stream pb.GpmService_InstallServiceStream) error {
	return s.manager.Install(ctx, &simpleInstallStream{stream: stream})
}

func (s *GpmServer) ListServiceVersions(ctx context.Context, req *pb.ListServiceVersionsReq, rsp *pb.ListServiceVersionsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Versions, err = s.manager.ListVersions(ctx, req.Name)
	return
}

func (s *GpmServer) UpgradeService(ctx context.Context, stream pb.GpmService_UpgradeServiceStream) error {
	return s.manager.Upgrade(ctx, &simpleUpgradeStream{stream: stream})
}

func (s *GpmServer) RollBackService(ctx context.Context, req *pb.RollbackServiceReq, rsp *pb.RollbackServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	err = s.manager.Rollback(ctx, req.Name, req.Revision)
	return
}

func (s *GpmServer) ForgetService(ctx context.Context, req *pb.ForgetServiceReq, rsp *pb.ForgetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	err = s.manager.Forget(ctx, req.Name, req.Revision)
	return
}

func (s *GpmServer) Ls(ctx context.Context, req *pb.LsReq, rsp *pb.LsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	rsp.Files, err = s.ftp.List(ctx, req.Path)
	return
}

func (s *GpmServer) Pull(ctx context.Context, req *pb.PullReq, stream pb.GpmService_PullStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	return s.ftp.Pull(ctx, req.Name, req.Dir, &simplePullSender{stream: stream})
}

func (s *GpmServer) Push(ctx context.Context, stream pb.GpmService_PushStream) (err error) {
	return s.ftp.Push(ctx, &simplePushReader{stream: stream})
}

func (s *GpmServer) Exec(ctx context.Context, req *pb.ExecReq, rsp *pb.ExecRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Result, err = s.ftp.Exec(ctx, req.In)
	return
}

func (s *GpmServer) Terminal(ctx context.Context, stream pb.GpmService_TerminalStream) error {
	return s.ftp.Terminal(ctx, &simpleTerminalStream{stream: stream})
}

type simpleUpdateSelfStream struct {
	stream pb.GpmService_UpdateSelfStream
}

func (s *simpleUpdateSelfStream) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simpleUpdateSelfStream) Send(msg interface{}) error {
	return s.stream.Send(&pb.UpdateSelfRsp{Result: msg.(*gpmv1.UpdateResult)})
}

func (s *simpleUpdateSelfStream) Close() error {
	return s.stream.Close()
}

type simpleWatchLogSender struct {
	stream pb.GpmService_WatchServiceLogStream
}

func (s *simpleWatchLogSender) Send(msg interface{}) error {
	return s.stream.Send(&pb.WatchServiceLogRsp{Log: msg.(*gpmv1.ServiceLog)})
}

func (s *simpleWatchLogSender) Close() error {
	return s.stream.Close()
}

type simpleInstallStream struct {
	stream pb.GpmService_InstallServiceStream
}

func (s *simpleInstallStream) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simpleInstallStream) Send(msg interface{}) error {
	return s.stream.Send(&pb.InstallServiceRsp{Result: msg.(*gpmv1.InstallServiceResult)})
}

func (s *simpleInstallStream) Close() error {
	return s.stream.Close()
}

type simpleUpgradeStream struct {
	stream pb.GpmService_UpgradeServiceStream
}

func (s *simpleUpgradeStream) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simpleUpgradeStream) Send(msg interface{}) error {
	return s.stream.Send(&pb.UpgradeServiceRsp{Result: msg.(*gpmv1.UpgradeServiceResult)})
}

func (s *simpleUpgradeStream) Close() error {
	return s.stream.Close()
}

type simplePullSender struct {
	stream pb.GpmService_PullStream
}

func (s *simplePullSender) Send(msg interface{}) error {
	return s.stream.Send(&pb.PullRsp{Result: msg.(*gpmv1.PullResult)})
}

func (s *simplePullSender) Close() error {
	return s.stream.Close()
}

type simplePushReader struct {
	stream pb.GpmService_PushStream
}

func (s *simplePushReader) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simplePushReader) Close() error {
	if e := s.stream.Send(&pb.PushRsp{}); e != nil {
		return e
	}
	return s.stream.Close()
}

type simpleTerminalStream struct {
	stream pb.GpmService_TerminalStream
}

func (s *simpleTerminalStream) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simpleTerminalStream) Send(msg interface{}) error {
	return s.stream.Send(&pb.TerminalRsp{Result: msg.(*gpmv1.TerminalResult)})
}

func (s *simpleTerminalStream) Close() error {
	return s.stream.Close()
}
