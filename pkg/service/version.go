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
	"runtime"
	"sort"
	"strings"
	"time"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/vine/lib/config"
	verrs "github.com/vine-io/vine/lib/errors"
	log "github.com/vine-io/vine/lib/logger"
)

func (g *manager) Install(ctx context.Context, stream IOStream) error {

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
		hname := hdr.Name
		if spec.HeaderTrimPrefix != "" {
			hname = strings.TrimPrefix(hname, spec.HeaderTrimPrefix)
		}
		fname := filepath.Join(dir, hname)
		if hdr.FileInfo().IsDir() {
			_ = os.MkdirAll(fname, os.ModePerm)
		} else {
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

	spec.InstallFlag = 1
	_, err = g.Create(ctx, spec)
	if err != nil {
		return err
	}

	log.Infof("install service %s@%s", spec.Name, spec.Version)

	return stream.Send(&gpmv1.InstallServiceResult{IsOk: true})
}

func (g *manager) ListVersions(ctx context.Context, name string) ([]*gpmv1.ServiceVersion, error) {
	vs, err := g.DB.ListServiceVersion(ctx, name)
	if err != nil {
		return nil, err
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Timestamp < vs[j].Timestamp
	})

	return vs, nil
}

func (g *manager) Upgrade(ctx context.Context, stream IOStream) error {

	var (
		file    *os.File
		err     error
		dst     string
		spec    *gpmv1.UpgradeSpec
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
		spec = b.Spec
		pack := b.Pack

		if file == nil {
			if err = b.Validate(); err != nil {
				return verrs.BadRequest(g.Name(), err.Error())
			}

			service, err = g.getService(ctx, spec.Name)
			if err != nil {
				return err
			}

			if service.Version == spec.Version {
				return verrs.Conflict(g.Name(), "version %s already exists", spec.Version)
			}

			_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "packages", spec.Name), 0o755)
			dst = filepath.Join(config.Get("root").String(""), "packages", spec.Name, spec.Name+"-"+spec.Version+".tar.gz")
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
	root := dir + "_" + spec.Version
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

		hname := hdr.Name
		if spec.HeaderTrimPrefix != "" {
			hname = strings.TrimPrefix(hname, spec.HeaderTrimPrefix)
		}
		fname := filepath.Join(dir, hname)
		if hdr.FileInfo().IsDir() {
			_ = os.MkdirAll(fname, os.ModePerm)
		} else {
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

	vf := spec.Version + "@" + time.Now().Format("20060102150405")
	log.Infof("service %s append version %s", service.Name, spec.Version)
	_ = ioutil.WriteFile(filepath.Join(config.Get("root").String(""), "services", spec.Name, "versions", vf), []byte(""), 0o777)

	service.Version = spec.Version
	g.DB.UpdateService(ctx, service)

	p = NewProcess(service)
	if isRunning {
		log.Infof("start service %s", service.Name)
		g.startService(ctx, p)
	}

	return stream.Send(&gpmv1.UpgradeServiceResult{IsOk: true})

}

func (g *manager) Rollback(ctx context.Context, name string, version string) error {
	s, err := g.getService(ctx, name)
	if err != nil {
		return err
	}

	vv, err := g.ListVersions(ctx, name)
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

func (g *manager) Forget(ctx context.Context, name string, version string) error {
	s, err := g.getService(ctx, name)
	if err != nil {
		return err
	}

	vv, err := g.ListVersions(ctx, name)
	if err != nil {
		return err
	}

	var v *gpmv1.ServiceVersion
	for _, item := range vv {
		if item.Version == version {
			v = item
		}
	}
	if v == nil {
		return verrs.NotFound(g.Name(), "invalid version '%s' of service:%s", version, name)
	}

	vf := version + "@" + time.Unix(v.Timestamp, 0).Format("20060102150405")
	vstore := filepath.Join(config.Get("root").String(""), "services", name, "versions", vf)
	log.Infof("remove %s@%s version file %s", name, version, vstore)
	if err = os.Remove(vstore); err != nil {
		log.Errorf("remove %s@%s version file: %v", name, version, err)
	}

	pkg := filepath.Join(config.Get("root").String(""), "packages", name, name+"-"+version+".tar.gz")
	log.Infof("remove %s@%s version package %s", name, version, pkg)
	if err = os.Remove(pkg); err != nil {
		log.Errorf("remove %s:%s version directory: %v", name, version, err)
	}

	sp := s.Dir + "_" + version
	log.Infof("remove %s@%s version directory %s", name, version, sp)
	if err = os.RemoveAll(sp); err != nil {
		log.Errorf("remove %s@%s version package: %v", name, version, err)
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	lIndex := 0
	switch runtime.GOOS {
	case "windows":
		lIndex = strings.LastIndex(name, "\\")
	default:
		lIndex = strings.LastIndex(name, "/")
	}
	if lIndex == -1 {
		lIndex = len(name)
	}
	err := os.MkdirAll(string([]rune(name)[0:lIndex]), 0o755)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
}
