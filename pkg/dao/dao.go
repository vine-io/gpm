package dao

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

	"github.com/gpm2/gpm/pkg/runtime/config"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
)

var (
	ErrTimeout  = errors.New("db request resource timeout")
	ErrNotFound = errors.New("resource not found")
)

type DB struct {
	Cfg *config.Config `inject:""`
}

func (db *DB) FindAllServices(ctx context.Context) ([]*gpmv1.Service, error) {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
		outs = make([]*gpmv1.Service, 0)
	)

	go func() {
		err := filepath.WalkDir(filepath.Join(db.Cfg.Root, "services"), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			sf := filepath.Join(path, d.Name()+".json")
			stat, _ := os.Stat(sf)
			if stat == nil {
				return nil
			}

			b, err := ioutil.ReadFile(sf)
			if err != nil {
				return nil
			}
			s := new(gpmv1.Service)
			err = json.Unmarshal(b, &s)
			if err != nil {
				return nil
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
	case <-ctx.Done():
		return nil, ErrTimeout
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
		f := filepath.Join(db.Cfg.Root, "services", name, name+".json")
		stat, _ := os.Stat(f)
		if stat == nil {
			ech <- fmt.Errorf("%w: service '%s'", ErrNotFound, name)
			return
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			ech <- err
			return
		}
		if err = json.Unmarshal(b, &out); err != nil {
			ech <- err
			return
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
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
		root := filepath.Join(db.Cfg.Root, "services", name, "versions")
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
				t, _ := time.Parse("20060102150405", parts[1])
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
	case <-ctx.Done():
		return nil, ErrTimeout
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
		_ = os.MkdirAll(filepath.Join(db.Cfg.Root, "services", s.Name), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(db.Cfg.Root, "logs", s.Name), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(db.Cfg.Root, "services", s.Name, "versions"), os.ModePerm)
		version := s.Version + "@" + time.Now().Format("20060102150405")
		_ = ioutil.WriteFile(filepath.Join(db.Cfg.Root, "services", s.Name, "versions", version), []byte(""), os.ModePerm)

		b, err := json.Marshal(s)
		if err != nil {
			ech <- err
			return
		}
		f := filepath.Join(db.Cfg.Root, "services", s.Name, s.Name+".json")
		if err = ioutil.WriteFile(f, b, os.ModePerm); err != nil {
			ech <- err
			return
		}

		out = s
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
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
		b, err := json.Marshal(s)
		if err != nil {
			ech <- err
			return
		}
		f := filepath.Join(db.Cfg.Root, "services", s.Name, s.Name+".json")
		if err = ioutil.WriteFile(f, b, os.ModePerm); err != nil {
			ech <- err
			return
		}

		out = s
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case e := <-ech:
		return nil, e
	case <-done:
		return out, nil
	}
}

func (db *DB) DeleteService(ctx context.Context, name string) error {
	var (
		done = make(chan struct{}, 1)
		ech  = make(chan error, 1)
	)

	go func() {
		_ = os.RemoveAll(filepath.Join(db.Cfg.Root, "services", name))
		_ = os.RemoveAll(filepath.Join(db.Cfg.Root, "logs", name))
		_ = os.RemoveAll(filepath.Join(db.Cfg.Root, "packages", name))

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrTimeout
	case e := <-ech:
		return e
	case <-done:
		return nil
	}
}
