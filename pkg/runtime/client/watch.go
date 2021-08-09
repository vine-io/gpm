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

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
)

type UpdateStream struct {
	s       pb.GpmService_UpdateSelfService
	version string
}

func NewUpdateStream(s pb.GpmService_UpdateSelfService) *UpdateStream {
	return &UpdateStream{s: s}
}

func (s *UpdateStream) Context() context.Context {
	return s.s.Context()
}

func (s *UpdateStream) Send(in *gpmv1.UpdateIn) error {
	err := s.s.Send(&pb.UpdateSelfReq{
		In: in,
	})
	return err
}

func (s *UpdateStream) Recv() (*gpmv1.UpdateResult, error) {
	rsp, err := s.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (s *UpdateStream) Close() error {
	return s.s.Close()
}

type ServiceLogWatcher struct {
	s pb.GpmService_WatchServiceLogService
}

func (w *ServiceLogWatcher) Context() context.Context {
	return w.s.Context()
}

func (w *ServiceLogWatcher) Next() (*gpmv1.ServiceLog, error) {
	rsp, err := w.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Log, nil
}

func (w *ServiceLogWatcher) Close() error {
	return w.s.Close()
}

type InstallStream struct {
	s pb.GpmService_InstallServiceService

	spec *gpmv1.ServiceSpec
}

func NewInstallStream(s pb.GpmService_InstallServiceService, spec *gpmv1.ServiceSpec) *InstallStream {
	return &InstallStream{s: s, spec: spec}
}

func (w *InstallStream) Context() context.Context {
	return w.s.Context()
}

func (w *InstallStream) Send(pack *gpmv1.Package) error {
	err := w.s.Send(&pb.InstallServiceReq{
		Spec: w.spec,
		Pack: pack,
	})
	return err
}

func (w *InstallStream) Recv() (*gpmv1.InstallServiceResult, error) {
	rsp, err := w.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (w *InstallStream) Close() error {
	return w.s.Close()
}

type UpgradeStream struct {
	s pb.GpmService_UpgradeServiceService

	name, version string
}

func NewUpgradeStream(s pb.GpmService_UpgradeServiceService, name, version string) *UpgradeStream {
	return &UpgradeStream{s: s, name: name, version: version}
}

func (s *UpgradeStream) Context() context.Context {
	return s.s.Context()
}

func (s *UpgradeStream) Send(pack *gpmv1.Package) error {
	err := s.s.Send(&pb.UpgradeServiceReq{
		Name:    s.name,
		Version: s.version,
		Pack:    pack,
	})
	return err
}

func (s *UpgradeStream) Recv() (*gpmv1.UpgradeServiceResult, error) {
	rsp, err := s.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (s *UpgradeStream) Close() error {
	return s.s.Close()
}

type PullWatcher struct {
	s pb.GpmService_PullService
}

func (w *PullWatcher) Context() context.Context {
	return w.s.Context()
}

func (w *PullWatcher) Next() (*gpmv1.PullResult, error) {
	rsp, err := w.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (w *PullWatcher) Close() error {
	return w.s.Close()
}

type PushStream struct {
	s pb.GpmService_PushService

	name, version string
}

func NewPushStream(s pb.GpmService_PushService) *PushStream {
	return &PushStream{s: s}
}

func (s *PushStream) Context() context.Context {
	return s.s.Context()
}

func (s *PushStream) Send(in *gpmv1.PushIn) error {
	err := s.s.Send(&pb.PushReq{
		In: in,
	})
	return err
}

func (s *PushStream) Recv() (*gpmv1.PushResult, error) {
	rsp, err := s.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (s *PushStream) Close() error {
	return s.s.Close()
}

type TerminalStream struct {
	s pb.GpmService_TerminalService

	name, version string
}

func NewTerminalStream(s pb.GpmService_TerminalService) *TerminalStream {
	return &TerminalStream{s: s}
}

func (s *TerminalStream) Context() context.Context {
	return s.s.Context()
}

func (s *TerminalStream) Send(in *gpmv1.TerminalIn) error {
	err := s.s.Send(&pb.TerminalReq{In: in})
	return err
}

func (s *TerminalStream) Recv() (*gpmv1.TerminalResult, error) {
	rsp, err := s.s.Recv()
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (s *TerminalStream) Close() error {
	return s.s.Close()
}
