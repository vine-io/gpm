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

package app

import (
	"context"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/domain"
	"github.com/vine-io/gpm/pkg/runtime/inject"
)

func init() {
	inject.ProvidePanic(new(gpmApp))
}

var _ GpmApp = (*gpmApp)(nil)

type gpmApp struct {
	G domain.Manager `inject:""`
	T domain.FTP     `inject:""`
}

func (g *gpmApp) Init() error {
	return g.G.Init()
}

func (g *gpmApp) Info(ctx context.Context) (*gpmv1.GpmInfo, error) {
	return g.G.Info(ctx)
}

func (g *gpmApp) UpdateSelf(ctx context.Context, stream domain.IOStream) error {
	return g.G.Update(ctx, stream)
}

func (g *gpmApp) ListService(ctx context.Context) ([]*gpmv1.Service, int64, error) {
	return g.G.List(ctx)
}

func (g *gpmApp) GetService(ctx context.Context, name string) (*gpmv1.Service, error) {
	return g.G.Get(ctx, name)
}

func (g *gpmApp) CreateService(ctx context.Context, spec *gpmv1.ServiceSpec) (*gpmv1.Service, error) {
	return g.G.Create(ctx, spec)
}

func (g *gpmApp) EditService(ctx context.Context, name string, spec *gpmv1.EditServiceSpec) (*gpmv1.Service, error) {
	return g.G.Edit(ctx, name, spec)
}

func (g *gpmApp) StartService(ctx context.Context, name string) (*gpmv1.Service, error) {
	return g.G.Start(ctx, name)
}

func (g *gpmApp) StopService(ctx context.Context, name string) (*gpmv1.Service, error) {
	return g.G.Stop(ctx, name)
}

func (g *gpmApp) RebootService(ctx context.Context, name string) (*gpmv1.Service, error) {
	return g.G.Reboot(ctx, name)
}

func (g *gpmApp) DeleteService(ctx context.Context, name string) (*gpmv1.Service, error) {
	return g.G.Delete(ctx, name)
}

func (g *gpmApp) WatchServiceLog(ctx context.Context, name string, number int64, follow bool, sender domain.IOWriter) error {
	return g.G.TailLog(ctx, name, number, follow, sender)
}

func (g *gpmApp) InstallService(ctx context.Context, stream domain.IOStream) error {
	return g.G.Install(ctx, stream)
}

func (g *gpmApp) ListServiceVersions(ctx context.Context, name string) ([]*gpmv1.ServiceVersion, error) {
	return g.G.ListVersions(ctx, name)
}

func (g *gpmApp) UpgradeService(ctx context.Context, stream domain.IOStream) error {
	return g.G.Upgrade(ctx, stream)
}

func (g *gpmApp) RollbackService(ctx context.Context, name string, version string) error {
	return g.G.Rollback(ctx, name, version)
}

func (g *gpmApp) Ls(ctx context.Context, root string) ([]*gpmv1.FileInfo, error) {
	return g.T.List(ctx, root)
}

func (g *gpmApp) Pull(ctx context.Context, name string, isDir bool, sender domain.IOWriter) error {
	return g.T.Pull(ctx, name, isDir, sender)
}

func (g *gpmApp) Push(ctx context.Context, stream domain.IOReader) error {
	return g.T.Push(ctx, stream)
}

func (g *gpmApp) Exec(ctx context.Context, in *gpmv1.ExecIn) (*gpmv1.ExecResult, error) {
	return g.T.Exec(ctx, in)
}

func (g *gpmApp) Terminal(ctx context.Context, stream domain.IOStream) error {
	return g.T.Terminal(ctx, stream)
}
