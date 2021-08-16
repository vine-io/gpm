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

	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
	"github.com/vine-io/vine/lib/config"
	log "github.com/vine-io/vine/lib/logger"
	verrs "github.com/vine-io/vine/proto/apis/errors"
)

func (g *gpm) InstallService(ctx context.Context, stream IOStream) error {

	var (
		file *os.File
		err  error
		dst  string
		spec *gpmv1.ServiceSpec
	)

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	var data interface{}
	for {
		data, err = stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}
		b := data.(*gpmv1.InstallServiceIn)
		spec = b.Spec
		pack := b.Pack

		if file == nil {
			if err = b.Validate(); err != nil {
				return verrs.BadRequest(g.Name(), err.Error())
			}

			v, _ := g.getService(ctx, spec.Name)
			if v != nil {
				return verrs.Conflict(g.Name(), "service '%s' already exists", spec.Name)
			}

			_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "packages", spec.Name), 0o755)
			dst = filepath.Join(config.Get("root").String(""), "packages", spec.Name, spec.Name+"-"+spec.Version+".tar.gz")
			file, err = os.Create(dst)
			if err != nil {
				return err
			}
		}

		if pack.Length > 0 {
			_, err = file.Write(pack.Chunk[0:pack.Length])
			if err != nil {
				return err
			}
		}
		if pack.IsOk {
			goto CHUNKED
		}

		if err == io.EOF {
			return nil
		}
	}

CHUNKED:
	_ = file.Close()
	file, err = os.Open(dst)
	if err != nil {
		return err
	}

	dir := spec.Dir
	root := dir + "_" + spec.Version
	_ = os.MkdirAll(root, 0o755)
	_ = os.Remove(dir)
	err = os.Symlink(root, dir)
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)
	for {
		hdr, e := tr.Next()
		if e != nil {
			if e == io.EOF {
				break
			} else {
				return e
			}
		}
		if hdr.FileInfo().IsDir() {
			_ = os.MkdirAll(filepath.Join(dir, hdr.Name), os.ModePerm)
		} else {
			fname := filepath.Join(dir, hdr.Name)
			f, e1 := createFile(fname)
			if e1 != nil {
				return e1
			}
			_, e1 = io.Copy(f, tr)
			if e1 != nil && e1 != io.EOF {
				f.Close()
				return e1
			}
			f.Close()
		}
	}

	_, err = g.CreateService(ctx, spec)
	if err != nil {
		return err
	}

	log.Infof("install service %s@%s", spec.Name, spec.Version)

	return stream.Send(&gpmv1.InstallServiceResult{IsOk: true})
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

func (g *gpm) UpgradeService(ctx context.Context, stream IOStream) error {

	var (
		file    *os.File
		err     error
		dst     string
		name    string
		version string
		service *gpmv1.Service
	)

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	var data interface{}
	for {
		data, err = stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}
		b := data.(*gpmv1.UpgradeServiceIn)
		name, version = b.Name, b.Version
		pack := b.Pack

		if file == nil {
			if err = b.Validate(); err != nil {
				return verrs.BadRequest(g.Name(), err.Error())
			}

			service, err = g.getService(ctx, name)
			if err != nil {
				return err
			}

			if service.Version == version {
				return verrs.Conflict(g.Name(), "version %s already exists", version)
			}

			_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "packages", name), 0o755)
			dst = filepath.Join(config.Get("root").String(""), "packages", name, name+"-"+version+".tar.gz")
			log.Infof("save package: %v", dst)
			file, err = os.Create(dst)
			if err != nil {
				//outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
				return err
			}
			//defer f.Close()
		}

		if pack.Length > 0 {
			_, err = file.Write(pack.Chunk[0:pack.Length])
			if err != nil {
				//outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
				return err
			}
		}

		if err == io.EOF {
			return nil
		}

		if pack.IsOk {
			goto CHUNKED
		}
	}

CHUNKED:

	g.RLock()
	p := g.ps[service.Name]
	g.RUnlock()

	isRunning := service.Status == gpmv1.StatusRunning
	if isRunning {
		log.Infof("stop service: %s", service.Name)
		g.stopService(ctx, p)
	}

	_ = file.Close()
	file, err = os.Open(dst)
	if err != nil {
		//outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
		return err
	}

	dir := service.Dir
	root := dir + "_" + version
	_ = os.MkdirAll(root, 0o755)
	_ = os.Remove(dir)
	log.Infof("relink %s -> %s", dir, root)
	err = os.Symlink(root, dir)
	if err != nil {
		//outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
		return err
	}

	log.Infof("unpack service %s package", service.Name)
	gr, err := gzip.NewReader(file)
	if err != nil {
		//outs <- &gpmv1.UpgradeServiceResult{Error: err.Error()}
		return err
	}
	tr := tar.NewReader(gr)
	for {
		hdr, e := tr.Next()
		if e != nil {
			if e == io.EOF {
				break
			} else {
				//outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
				return e
			}
		}
		fname := filepath.Join(dir, hdr.Name)
		f, e1 := createFile(fname)
		if e1 != nil {
			//outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
			return e1
		}
		_, e1 = io.Copy(f, tr)
		if e1 != nil && e1 != io.EOF {
			//outs <- &gpmv1.UpgradeServiceResult{Error: e.Error()}
			f.Close()
			return e1
		}
		f.Close()
	}

	vf := version + "@" + time.Now().Format("20060102150405")
	log.Infof("service %s append version %s", service.Name, version)
	_ = ioutil.WriteFile(filepath.Join(config.Get("root").String(""), "services", name, "versions", vf), []byte(""), 0o777)

	service.Version = version
	g.DB.UpdateService(ctx, service)

	p = NewProcess(service)
	if isRunning {
		log.Infof("start service %s", service.Name)
		g.startService(ctx, p)
	}

	return stream.Send(&gpmv1.UpgradeServiceResult{IsOk: true})

}

func (g *gpm) RollbackService(ctx context.Context, name string, version string) error {
	s, err := g.getService(ctx, name)
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

	g.RLock()
	p := g.ps[s.Name]
	g.RUnlock()
	isRunning := s.Status == gpmv1.StatusRunning
	if isRunning {
		g.stopService(ctx, p)
	}
	dir := s.Dir
	root := s.Dir + "_" + version
	log.Infof("relink %s -> %s", dir, root)
	_ = os.Remove(dir)
	err = os.Symlink(root, dir)
	if err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}

	s.Version = version
	_, err = g.DB.UpdateService(ctx, s)
	if err != nil {
		return err
	}

	p = NewProcess(s)
	if isRunning {
		g.startService(ctx, p)
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0o755)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
}
