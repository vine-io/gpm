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

package interfaces

import (
	"context"

	pb "github.com/vine-io/gpm/api/service/gpm/v1"
	verrs "github.com/vine-io/vine/lib/errors"
)

func (s *GpmAPI) Healthz(ctx context.Context, _ *pb.Empty, rsp *pb.Empty) error {
	return nil
}

func (s *GpmAPI) UpdateSelf(ctx context.Context, stream pb.GpmService_UpdateSelfStream) error {
	return s.H.UpdateSelf(ctx, &simpleUpdateSelfStream{stream: stream})
}

func (s *GpmAPI) Info(ctx context.Context, _ *pb.InfoReq, rsp *pb.InfoRsp) (err error) {
	rsp.Gpm, err = s.H.Info(ctx)
	return
}

func (s *GpmAPI) ListService(ctx context.Context, req *pb.ListServiceReq, rsp *pb.ListServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Services, rsp.Total, err = s.H.ListService(ctx)
	return
}

func (s *GpmAPI) GetService(ctx context.Context, req *pb.GetServiceReq, rsp *pb.GetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.GetService(ctx, req.Name)
	return
}

func (s *GpmAPI) CreateService(ctx context.Context, req *pb.CreateServiceReq, rsp *pb.CreateServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.CreateService(ctx, req.Spec)
	return
}

func (s *GpmAPI) EditService(ctx context.Context, req *pb.EditServiceReq, rsp *pb.EditServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.EditService(ctx, req.Name, req.Spec)
	return
}

func (s *GpmAPI) StartService(ctx context.Context, req *pb.StartServiceReq, rsp *pb.StartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StartService(ctx, req.Name)
	return
}

func (s *GpmAPI) StopService(ctx context.Context, req *pb.StopServiceReq, rsp *pb.StopServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StopService(ctx, req.Name)
	return
}

func (s *GpmAPI) RebootService(ctx context.Context, req *pb.RebootServiceReq, rsp *pb.RebootServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.RebootService(ctx, req.Name)
	return
}

func (s *GpmAPI) DeleteService(ctx context.Context, req *pb.DeleteServiceReq, rsp *pb.DeleteServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.DeleteService(ctx, req.Name)
	return
}

func (s *GpmAPI) WatchServiceLog(ctx context.Context, req *pb.WatchServiceLogReq, stream pb.GpmService_WatchServiceLogStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	return s.H.WatchServiceLog(ctx, req.Name, req.Number, req.Follow, &simpleWatchLogSender{stream: stream})
}

func (s *GpmAPI) InstallService(ctx context.Context, stream pb.GpmService_InstallServiceStream) error {
	return s.H.InstallService(ctx, &simpleInstallStream{stream: stream})
}

func (s *GpmAPI) ListServiceVersions(ctx context.Context, req *pb.ListServiceVersionsReq, rsp *pb.ListServiceVersionsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Versions, err = s.H.ListServiceVersions(ctx, req.Name)
	return
}

func (s *GpmAPI) UpgradeService(ctx context.Context, stream pb.GpmService_UpgradeServiceStream) error {
	return s.H.UpgradeService(ctx, &simpleUpgradeStream{stream: stream})
}

func (s *GpmAPI) RollBackService(ctx context.Context, req *pb.RollbackServiceReq, rsp *pb.RollbackServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	err = s.H.RollbackService(ctx, req.Name, req.Revision)
	return
}

func (s *GpmAPI) Ls(ctx context.Context, req *pb.LsReq, rsp *pb.LsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	rsp.Files, err = s.H.Ls(ctx, req.Path)
	return
}

func (s *GpmAPI) Pull(ctx context.Context, req *pb.PullReq, stream pb.GpmService_PullStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	return s.H.Pull(ctx, req.Name, req.Dir, &simplePullSender{stream: stream})
}

func (s *GpmAPI) Push(ctx context.Context, stream pb.GpmService_PushStream) (err error) {
	return s.H.Push(ctx, &simplePushReader{stream: stream})
}

func (s *GpmAPI) Exec(ctx context.Context, req *pb.ExecReq, rsp *pb.ExecRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Result, err = s.H.Exec(ctx, req.In)
	return
}

func (s *GpmAPI) Terminal(ctx context.Context, stream pb.GpmService_TerminalStream) error {
	return s.H.Terminal(ctx, &simpleTerminalStream{stream: stream})
}
