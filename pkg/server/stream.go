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
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
)

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

type simplePushStream struct {
	stream pb.GpmService_PushStream
}

func (s *simplePushStream) Recv() (interface{}, error) {
	b, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return b.In, nil
}

func (s *simplePushStream) Send(msg interface{}) error {
	return s.stream.Send(&pb.PushRsp{Result: msg.(*gpmv1.PushResult)})
}

func (s *simplePushStream) Close() error {
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