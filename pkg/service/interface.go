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

	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
)

type Gpm interface {
	Init() error
	Info(context.Context) (*gpmv1.GpmInfo, error)
	UpdateSelf(context.Context, IOStream) error
	ListService(context.Context) ([]*gpmv1.Service, int64, error)
	GetService(context.Context, string) (*gpmv1.Service, error)
	CreateService(context.Context, *gpmv1.ServiceSpec) (*gpmv1.Service, error)
	EditService(context.Context, string, *gpmv1.EditServiceSpec) (*gpmv1.Service, error)
	StartService(context.Context, string) (*gpmv1.Service, error)
	StopService(context.Context, string) (*gpmv1.Service, error)
	RebootService(context.Context, string) (*gpmv1.Service, error)
	DeleteService(context.Context, string) (*gpmv1.Service, error)
	WatchServiceLog(context.Context, string, int64, bool, IOWriter) error

	InstallService(context.Context, IOStream) error
	ListServiceVersions(context.Context, string) ([]*gpmv1.ServiceVersion, error)
	UpgradeService(context.Context, IOStream) error
	RollbackService(context.Context, string, string) error

	Ls(context.Context, string) ([]*gpmv1.FileInfo, error)
	Pull(context.Context, string, bool, IOWriter) error
	Push(context.Context, IOReader) error
	Exec(context.Context, *gpmv1.ExecIn) (*gpmv1.ExecResult, error)
	Terminal(context.Context, IOStream) error
}

type IOWriter interface {
	Send(interface{}) error
	Close() error
}

type IOReader interface {
	Recv() (interface{}, error)
	Close() error
}

type IOStream interface {
	Recv() (interface{}, error)
	Send(msg interface{}) error
	Close() error
}
