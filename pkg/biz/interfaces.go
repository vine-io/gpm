package biz

import (
	"context"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
)

type Manager interface {
	Init() error
	Info(context.Context) (*gpmv1.GpmInfo, error)
	Update(context.Context, IOStream) error
	List(context.Context) ([]*gpmv1.Service, int64, error)
	Get(context.Context, string) (*gpmv1.Service, error)
	Create(context.Context, *gpmv1.ServiceSpec) (*gpmv1.Service, error)
	Edit(context.Context, string, *gpmv1.EditServiceSpec) (*gpmv1.Service, error)
	Start(context.Context, string) (*gpmv1.Service, error)
	Stop(context.Context, string) (*gpmv1.Service, error)
	Restart(context.Context, string) (*gpmv1.Service, error)
	Delete(context.Context, string) (*gpmv1.Service, error)
	TailLog(context.Context, string, int64, bool, IOWriter) error

	Install(context.Context, IOStream) error
	ListVersions(context.Context, string) ([]*gpmv1.ServiceVersion, error)
	Upgrade(context.Context, IOStream) error
	Rollback(context.Context, string, string) error
	Forget(context.Context, string, string) error
}

type FTP interface {
	List(context.Context, string) ([]*gpmv1.FileInfo, error)
	Pull(context.Context, string, bool, IOWriter) error
	Push(context.Context, IOReader) error
	Exec(context.Context, *gpmv1.ExecIn) (*gpmv1.ExecResult, error)
	Terminal(context.Context, IOStream) error
}
