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

package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"github.com/vine-io/vine/lib/config"
	"gopkg.in/yaml.v3"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
)

var (
	ErrTimeout  = errors.New("db request resource timeout")
	ErrNotFound = errors.New("resource not found")
)

type DB struct{}

func (db *DB) FindAllServices(ctx context.Context) ([]*gpmv1.Service, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		outs = make([]*gpmv1.Service, 0)
	)

	go func() {
		err := filepath.WalkDir(filepath.Join(config.Get("root").String(""), "services"), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			sf := filepath.Join(path, d.Name()+".yml")
			stat, _ := os.Stat(sf)
			if stat == nil {
				sf = filepath.Join(path, d.Name()+".json")
				stat, _ = os.Stat(sf)
				if stat == nil {
					return nil
				}
			}

			b, err := ioutil.ReadFile(sf)
			if err != nil {
				return nil
			}
			s := new(gpmv1.Service)

			switch filepath.Ext(sf) {
			case ".yml":
				err = yaml.Unmarshal(b, &s)
				if err != nil {
					return nil
				}
			case ".json":
				err = json.Unmarshal(b, &s)
				if err != nil {
					return nil
				}
			}

			outs = append(outs, s)

			return nil
		})
		if err != nil {
			ech <- err
		} else {
			done <- struct{}{}
		}
	}()

	select {
	case e := <-ech:
		return nil, e
	case <-done:
		return outs, nil
	}
}

func (db *DB) FindService(ctx context.Context, name string) (*gpmv1.Service, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		out  = new(gpmv1.Service)
	)

	go func() {
		f := filepath.Join(config.Get("root").String(""), "services", name, name+".yml")
		stat, _ := os.Stat(f)
		if stat == nil {
			f = filepath.Join(config.Get("root").String(""), "services", name, name+".json")
			stat, _ = os.Stat(f)
			if stat == nil {
				ech <- fmt.Errorf("%w: service '%s'", ErrNotFound, name)
				return
			}
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			ech <- err
			return
		}

		switch filepath.Ext(f) {
		case ".json":
			if err = json.Unmarshal(b, &out); err != nil {
				ech <- err
				return
			}
		case ".yml":
			if err = yaml.Unmarshal(b, &out); err != nil {
				ech <- err
				return
			}
		}

		done <- struct{}{}
	}()

	select {
	case e := <-ech:
		return nil, e
	case <-done:
		return out, nil
	}
}

func (db *DB) ListServiceVersion(ctx context.Context, name string) ([]*gpmv1.ServiceVersion, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		outs = make([]*gpmv1.ServiceVersion, 0)
	)

	go func() {
		root := filepath.Join(config.Get("root").String(""), "services", name, "versions")
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() && path != root {
				return filepath.SkipDir
			}

			if !strings.Contains(d.Name(), "@") {
				return nil
			}

			parts := strings.Split(d.Name(), "@")
			if len(parts) > 1 {
				loc, _ := time.LoadLocation("Local")
				var t time.Time
				if loc != nil {
					t, _ = time.ParseInLocation("20060102150405", parts[1], loc)
				} else {
					t, _ = time.Parse("20060102150405", parts[1])
				}
				sv := &gpmv1.ServiceVersion{
					Name:      name,
					Version:   parts[0],
					Timestamp: t.Unix(),
				}
				outs = append(outs, sv)
			}

			return nil
		})
		if err != nil {
			ech <- err
		} else {
			done <- struct{}{}
		}
	}()

	select {
	case e := <-ech:
		return nil, e
	case <-done:
		return outs, nil
	}
}

func (db *DB) CreateService(ctx context.Context, s *gpmv1.Service) (*gpmv1.Service, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		out  = new(gpmv1.Service)
	)

	go func() {
		_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "services", s.Name), 0o777)
		_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "logs", s.Name), 0o777)
		_ = os.MkdirAll(filepath.Join(config.Get("root").String(""), "services", s.Name, "versions"), 0o777)
		version := s.Version + "@" + time.Now().Format("20060102150405")
		_ = ioutil.WriteFile(filepath.Join(config.Get("root").String(""), "services", s.Name, "versions", version), []byte(""), 0o777)

		b, err := yaml.Marshal(s)
		if err != nil {
			ech <- err
			return
		}
		f := filepath.Join(config.Get("root").String(""), "services", s.Name, s.Name+".yml")
		if err = ioutil.WriteFile(f, b, 0o777); err != nil {
			ech <- err
			return
		}

		out = s
		done <- struct{}{}
	}()

	select {
	case e := <-ech:
		return nil, e
	case <-done:
		return out, nil
	}
}

func (db *DB) UpdateService(ctx context.Context, s *gpmv1.Service) (*gpmv1.Service, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		out  = new(gpmv1.Service)
	)

	go func() {
		b, err := yaml.Marshal(s)
		if err != nil {
			ech <- err
			return
		}
		f := filepath.Join(config.Get("root").String(""), "services", s.Name, s.Name+".yml")
		if err = ioutil.WriteFile(f, b, 0o777); err != nil {
			ech <- err
			return
		}

		out = s
		done <- struct{}{}
	}()

	select {
	case e := <-ech:
		return nil, e
	case <-done:
		return out, nil
	}
}

func (db *DB) DeleteService(ctx context.Context, name string) error {
	var (
		done = make(chan struct{}, 1)
	)

	go func() {
		_ = os.RemoveAll(filepath.Join(config.Get("root").String(""), "services", name))
		_ = os.RemoveAll(filepath.Join(config.Get("root").String(""), "logs", name))
		_ = os.RemoveAll(filepath.Join(config.Get("root").String(""), "packages", name))

		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	}
}
