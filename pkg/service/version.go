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
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	log "github.com/lack-io/vine/lib/logger"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

func (g *gpm) InstallService(
	ctx context.Context,
	spec *gpmv1.ServiceSpec,
	in <-chan *gpmv1.Package,
) (<-chan *gpmv1.InstallServiceResult, error) {

	v, _ := g.GetService(ctx, spec.Name)
	if v != nil {
		return nil, verrs.Conflict(g.Name(), "service '%s' already exists", spec.Name)
	}

	outs := make(chan *gpmv1.InstallServiceResult, 10)

	go func() {
		_ = os.MkdirAll(filepath.Join(g.Cfg.Root, "packages"), 0755)
		pack := filepath.Join(g.Cfg.Root, "packages", spec.Name, spec.Name+"-"+spec.Version+".tar.gz")
		f, err := os.Create(pack)
		if err != nil {
			outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
			return
		}
		defer f.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case p, ok := <-in:
				if !ok {
					return
				}
				_, err = f.Write(p.Chunk)
				if err != nil {
					outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
					return
				}
				if p.IsOk {
					goto CHUNKED
				}
			}
		}

	CHUNKED:
		_ = f.Close()
		f, err = os.Open(pack)
		if err != nil {
			outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
			return
		}

		dir := spec.Dir
		root := dir + "_" + spec.Version
		_ = os.MkdirAll(root, 0755)
		_ = os.Remove(dir)
		err = os.Symlink(root, dir)
		if err != nil {
			outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
			return
		}

		gr, err := gzip.NewReader(f)
		if err != nil {
			outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
			return
		}
		tr := tar.NewReader(gr)
		for {
			hdr, e := tr.Next()
			if e != nil {
				if e == io.EOF {
					break
				} else {
					outs <- &gpmv1.InstallServiceResult{Error: e.Error()}
					return
				}
			}
			fname := filepath.Join(dir, hdr.Name)
			file, e1 := createFile(fname)
			if e1 != nil {
				outs <- &gpmv1.InstallServiceResult{Error: e.Error()}
				return
			}
			_, e1 = io.Copy(file, tr)
			if e1 != nil && e1 != io.EOF {
				outs <- &gpmv1.InstallServiceResult{Error: e.Error()}
				file.Close()
				return
			}
			file.Close()
		}

		_, err = g.CreateService(ctx, spec)
		if err != nil {
			outs <- &gpmv1.InstallServiceResult{Error: err.Error()}
			return
		}

		outs <- &gpmv1.InstallServiceResult{IsOk: true}
		log.Infof("install service %s@%s", spec.Name, spec.Version)
		return
	}()

	return outs, nil
}

func (g *gpm) ListServiceVersions(ctx context.Context, name string) ([]*gpmv1.ServiceVersion, error) {
	vs, err := g.DB.ListServiceVersion(ctx, name)
	if err != nil {
		return nil, err
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Timestamp < vs[j].Timestamp
	})

	return vs, nil
}

func (g *gpm) UpgradeService(
	ctx context.Context, name, version string,
	in <-chan *gpmv1.Package,
) (<-chan *gpmv1.UpgradeServiceResult, error) {

	service, err := g.GetService(ctx, name)
	if err != nil {
		return nil, err
	}

	outs := make(chan *gpmv1.UpgradeServiceResult, 10)

	go func() {
		_ = os.MkdirAll(filepath.Join(g.Cfg.Root, "packages"), 0755)
		pack := filepath.Join(g.Cfg.Root, "packages", name, name+"-"+version+".tar.gz")
		f, err := os.Create(pack)
		if err != nil {
			outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
			return
		}
		defer f.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case p, ok := <-in:
				if !ok {
					return
				}
				_, err = f.Write(p.Chunk)
				if err != nil {
					outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
					return
				}
				if p.IsOk {
					goto CHUNKED
				}
			}
		}

	CHUNKED:

		isRunning := service.Status == gpmv1.StatusRunning
		if isRunning {
			g.StopService(ctx, service.Name)
		}

		_ = f.Close()
		f, err = os.Open(pack)
		if err != nil {
			outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
			return
		}

		dir := service.Dir
		root := dir + "_" + version
		_ = os.MkdirAll(root, 0755)
		_ = os.Remove(dir)
		err = os.Symlink(root, dir)
		if err != nil {
			outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
			return
		}

		gr, err := gzip.NewReader(f)
		if err != nil {
			outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
			return
		}
		tr := tar.NewReader(gr)
		for {
			hdr, e := tr.Next()
			if e != nil {
				if e == io.EOF {
					break
				} else {
					outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
					return
				}
			}
			fname := filepath.Join(dir, hdr.Name)
			file, e1 := createFile(fname)
			if e1 != nil {
				outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
				return
			}
			_, e1 = io.Copy(file, tr)
			if e1 != nil && e1 != io.EOF {
				outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
				file.Close()
				return
			}
			file.Close()
		}

		vf := service.Version + "@" + time.Now().Format("20060102150405")
		_ = ioutil.WriteFile(filepath.Join(g.Cfg.Root, "services", name, "versions", vf), []byte(""), os.ModePerm)

		service.Version = version
		g.DB.UpdateService(ctx, service)
		if isRunning {
			g.StartService(ctx, service.Name)
		}

		outs <- &gpmv1.UpgradeServiceResult{IsOk: true}
		return
	}()

	return outs, nil

}

func (g *gpm) RollbackService(ctx context.Context, name string, version string) error {
	s, err := g.GetService(ctx, name)
	if err != nil {
		return err
	}

	vv, err := g.ListServiceVersions(ctx, name)
	if err != nil {
		return err
	}

	exists := false
	for _, v := range vv {
		if v.Version == version {
			exists = true
		}
	}
	if !exists {
		return verrs.NotFound(g.Name(), "invalid version '%s' of service:%s", version, name)
	}

	isRunning := s.Status == gpmv1.StatusRunning
	if isRunning {
		g.StopService(ctx, name)
	}
	_ = os.Remove(s.Dir)
	err = os.Symlink(s.Dir+"_"+version, s.Dir)
	if err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}

	s.Version = version
	_, err = g.DB.UpdateService(ctx, s)
	if err != nil {
		return err
	}
	if isRunning {
		g.StartService(ctx, name)
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
