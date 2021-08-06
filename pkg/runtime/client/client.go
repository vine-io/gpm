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

package client

import (
	"context"

	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/ssl"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/vine/core/client"
	"github.com/lack-io/vine/core/client/grpc"
)

type SimpleClient struct {
	cc pb.GpmService

	addr string
}

func New() *SimpleClient {
	sc := &SimpleClient{}

	tls, _ := ssl.GetTLS()

	conn := grpc.NewClient(client.Retries(0), grpc.AuthTLS(tls))
	sc.cc = pb.NewGpmService(runtime.GpmName, conn)

	return sc
}

func (s *SimpleClient) Healthz(ctx context.Context, opts ...client.CallOption) error {
	_, err := s.cc.Healthz(ctx, &pb.Empty{}, opts...)
	return err
}

func (s *SimpleClient) Info(ctx context.Context, opts ...client.CallOption) (*gpmv1.GpmInfo, error) {
	rsp, err := s.cc.Info(ctx, &pb.InfoReq{}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Gpm, nil
}

func (s *SimpleClient) Update(ctx context.Context, opts ...client.CallOption) (*UpdateStream, error) {
	stream, err := s.cc.UpdateSelf(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewUpdateStream(stream), nil
}

func (s *SimpleClient) ListService(ctx context.Context, opts ...client.CallOption) ([]*gpmv1.Service, int64, error) {
	rsp, err := s.cc.ListService(ctx, &pb.ListServiceReq{}, opts...)
	if err != nil {
		return nil, 0, err
	}
	return rsp.Services, rsp.Total, nil
}

func (s *SimpleClient) GetService(ctx context.Context, name string, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.GetService(ctx, &pb.GetServiceReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) CreateService(ctx context.Context, spec *gpmv1.ServiceSpec, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.CreateService(ctx, &pb.CreateServiceReq{Spec: spec}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) EditService(ctx context.Context, name string, spec *gpmv1.EditServiceSpec, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.EditService(ctx, &pb.EditServiceReq{Name: name, Spec: spec}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) StartService(ctx context.Context, name string, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.StartService(ctx, &pb.StartServiceReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) StopService(ctx context.Context, name string, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.StopService(ctx, &pb.StopServiceReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) RebootService(ctx context.Context, name string, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.RebootService(ctx, &pb.RebootServiceReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) DeleteService(ctx context.Context, name string, opts ...client.CallOption) (*gpmv1.Service, error) {
	rsp, err := s.cc.DeleteService(ctx, &pb.DeleteServiceReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Service, nil
}

func (s *SimpleClient) WatchServiceLog(ctx context.Context, name string, n int64, f bool, opts ...client.CallOption) (*ServiceLogWatcher, error) {
	rsp, err := s.cc.WatchServiceLog(ctx, &pb.WatchServiceLogReq{
		Name:   name,
		Number: n,
		Follow: f,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &ServiceLogWatcher{s: rsp}, nil
}

func (s *SimpleClient) InstallService(ctx context.Context, spec *gpmv1.ServiceSpec, opts ...client.CallOption) (*InstallStream, error) {
	stream, err := s.cc.InstallService(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewInstallStream(stream, spec), nil
}

func (s *SimpleClient) ListServiceVersions(ctx context.Context, name string, opts ...client.CallOption) ([]*gpmv1.ServiceVersion, error) {
	rsp, err := s.cc.ListServiceVersions(ctx, &pb.ListServiceVersionsReq{Name: name}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Versions, nil
}

func (s *SimpleClient) UpgradeService(ctx context.Context, name, version string, opts ...client.CallOption) (*UpgradeStream, error) {
	stream, err := s.cc.UpgradeService(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewUpgradeStream(stream, name, version), nil
}

func (s *SimpleClient) RollBackService(ctx context.Context, name, revision string, opts ...client.CallOption) error {
	_, err := s.cc.RollBackService(ctx, &pb.RollbackServiceReq{Name: name, Revision: revision}, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (s *SimpleClient) Ls(ctx context.Context, path string, opts ...client.CallOption) ([]*gpmv1.FileInfo, error) {
	rsp, err := s.cc.Ls(ctx, &pb.LsReq{Path: path}, opts...)
	if err != nil {
		return nil, err
	}
	return rsp.Files, nil
}

func (s *SimpleClient) Pull(ctx context.Context, name string, isDir bool, opts ...client.CallOption) (*PullWatcher, error) {
	stream, err := s.cc.Pull(ctx, &pb.PullReq{Name: name, Dir: isDir}, opts...)
	if err != nil {
		return nil, err
	}
	return &PullWatcher{s: stream}, nil
}

func (s *SimpleClient) Push(ctx context.Context, opts ...client.CallOption) (*PushStream, error) {
	stream, err := s.cc.Push(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewPushStream(stream), nil
}

func (s *SimpleClient) Exec(ctx context.Context, in *gpmv1.ExecIn, opts ...client.CallOption) (*ExecWatcher, error) {
	rsp, err := s.cc.Exec(ctx, &pb.ExecReq{In: in}, opts...)
	if err != nil {
		return nil, err
	}
	return &ExecWatcher{s: rsp}, nil
}

func (s *SimpleClient) Terminal(ctx context.Context, opts ...client.CallOption) (*TerminalStream, error) {
	stream, err := s.cc.Terminal(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewTerminalStream(stream), nil
}
