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
	"io"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	log "github.com/lack-io/vine/lib/logger"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

func (s *server) Healthz(ctx context.Context, _ *pb.Empty, rsp *pb.Empty) error {
	return nil
}

func (s *server) UpdateSelf(ctx context.Context, stream pb.GpmService_UpdateSelfStream) (err error) {
	req, err := stream.Recv()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	in := make(chan *gpmv1.UpdateIn, 10)

	in <- req.In
	outs, e := s.H.UpdateSelf(ctx, req.In.Version, in)
	if e != nil {
		err = e
		return
	}

	go func() {
		defer close(in)
		for {
			req, err = stream.Recv()
			if err != nil {
				return
			}

			in <- req.In
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-outs:
			_ = stream.Send(&pb.UpdateSelfRsp{Result: out})
		}
	}
}

func (s *server) Info(ctx context.Context, _ *pb.InfoReq, rsp *pb.InfoRsp) (err error) {
	rsp.Gpm, err = s.H.Info(ctx)
	return
}

func (s *server) ListService(ctx context.Context, req *pb.ListServiceReq, rsp *pb.ListServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Services, rsp.Total, err = s.H.ListService(ctx)
	return
}

func (s *server) GetService(ctx context.Context, req *pb.GetServiceReq, rsp *pb.GetServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.GetService(ctx, req.Name)
	return
}

func (s *server) CreateService(ctx context.Context, req *pb.CreateServiceReq, rsp *pb.CreateServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.CreateService(ctx, req.Spec)
	return
}

func (s *server) EditService(ctx context.Context, req *pb.EditServiceReq, rsp *pb.EditServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.EditService(ctx, req.Name, req.Spec)
	return
}

func (s *server) StartService(ctx context.Context, req *pb.StartServiceReq, rsp *pb.StartServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StartService(ctx, req.Name)
	return
}

func (s *server) StopService(ctx context.Context, req *pb.StopServiceReq, rsp *pb.StopServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.StopService(ctx, req.Name)
	return
}

func (s *server) RebootService(ctx context.Context, req *pb.RebootServiceReq, rsp *pb.RebootServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.RebootService(ctx, req.Name)
	return
}

func (s *server) DeleteService(ctx context.Context, req *pb.DeleteServiceReq, rsp *pb.DeleteServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Service, err = s.H.DeleteService(ctx, req.Name)
	return
}

func (s *server) WatchServiceLog(ctx context.Context, req *pb.WatchServiceLogReq, stream pb.GpmService_WatchServiceLogStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	outs, e := s.H.WatchServiceLog(ctx, req.Name, req.Number, req.Follow)
	if e != nil {
		err = e
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-outs:
			if !ok {
				return
			}
			e = stream.Send(&pb.WatchServiceLogRsp{Log: out})
			if e != nil && e != io.EOF {
				err = e
				return
			}
		}
	}
}

func (s *server) InstallService(ctx context.Context, stream pb.GpmService_InstallServiceStream) (err error) {
	req, err := stream.Recv()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	in := make(chan *gpmv1.Package, 10)
	in <- req.Pack
	outs, e := s.H.InstallService(ctx, req.Spec, in)
	if e != nil {
		err = e
		return
	}

	go func() {
		defer close(in)
		for {
			req, err = stream.Recv()
			if err != nil {
				return
			}

			in <- req.Pack
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-outs:
			if !ok {
				return
			}
			rsp := &pb.InstallServiceRsp{Result: out}
			_ = stream.Send(rsp)
			if out.IsOk {
				err = nil
				return
			}
		}
	}
}

func (s *server) ListServiceVersions(ctx context.Context, req *pb.ListServiceVersionsReq, rsp *pb.ListServiceVersionsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Versions, err = s.H.ListServiceVersions(ctx, req.Name)
	return
}

func (s *server) UpgradeService(ctx context.Context, stream pb.GpmService_UpgradeServiceStream) (err error) {
	req, err := stream.Recv()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	in := make(chan *gpmv1.Package, 10)

	in <- req.Pack
	outs, e := s.H.UpgradeService(ctx, req.Name, req.Version, in)
	if e != nil {
		err = e
		return
	}

	go func() {
		defer close(in)
		for {
			req, err = stream.Recv()
			if err != nil {
				return
			}

			in <- req.Pack
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-outs:
			if out.Error != "" {
				return verrs.InternalServerError(s.Name(), out.Error)
			}
			rsp := &pb.UpgradeServiceRsp{Result: out}
			_ = stream.Send(rsp)
			if out.IsOk {
				err = nil
				return
			}
		}
	}
}

func (s *server) RollBackService(ctx context.Context, req *pb.RollbackServiceReq, rsp *pb.RollbackServiceRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	err = s.H.RollbackService(ctx, req.Name, req.Revision)
	return
}

func (s *server) Ls(ctx context.Context, req *pb.LsReq, rsp *pb.LsRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	rsp.Files, err = s.H.Ls(ctx, req.Path)
	return
}

func (s *server) Pull(ctx context.Context, req *pb.PullReq, stream pb.GpmService_PullStream) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	outs, e := s.H.Pull(ctx, req.Name, req.Dir)
	if e != nil {
		err = e
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-outs:
			rsp := &pb.PullRsp{Result: result}
			e = stream.Send(rsp)
			if e != nil && e != io.EOF {
				err = e
				return
			}
			if result.Finished {
				return
			}
		}
	}
}

func (s *server) Push(ctx context.Context, stream pb.GpmService_PushStream) (err error) {
	req, err := stream.Recv()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	in := make(chan *gpmv1.PushIn, 10)
	defer close(in)

	in <- req.In
	outs, e := s.H.Push(ctx, req.In.Dst, req.In.Name, in)
	if e != nil {
		err = e
		return
	}

	go func() {
		for {
			req, err = stream.Recv()
			if err != nil && err != io.EOF {
				log.Errorf("receive data: %v", err)
				return
			}

			in <- req.In

			if io.EOF != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-outs:
			if !ok {
				return
			}
			_ = stream.Send(&pb.PushRsp{Result: out})
		}
	}
}

func (s *server) Exec(ctx context.Context, req *pb.ExecReq, rsp *pb.ExecRsp) (err error) {
	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	rsp.Result, err = s.H.Exec(ctx, req.In)
	return
}

func (s *server) Terminal(ctx context.Context, stream pb.GpmService_TerminalStream) (err error) {
	req, err := stream.Recv()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	if err = req.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}

	in := make(chan *gpmv1.TerminalIn, 10)

	in <- req.In
	outs, e := s.H.Terminal(ctx, in)
	if e != nil {
		err = e
		return
	}

	go func() {
		defer close(in)
		for {
			req, err = stream.Recv()
			if err != nil {
				return
			}

			in <- req.In
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case out := <-outs:
			rsp := &pb.TerminalRsp{Result: out}
			e = stream.Send(rsp)
			if e != nil {
				err = e
				return
			}
			if out.IsOk {
				return
			}
		}
	}
}
