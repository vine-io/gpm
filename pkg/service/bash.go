package service

import (
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

func (g *gpm) Ls(ctx context.Context, root string) ([]*gpmv1.FileInfo, error) {
	stat, _ := os.Stat(root)
	if stat == nil {
		return nil, verrs.NotFound(g.Name(), "invalid path '%s'", root)
	}

	files := make([]*gpmv1.FileInfo, 0)
	if !stat.IsDir() {
		files = append(files, &gpmv1.FileInfo{
			Name:    stat.Name(),
			Size:    stat.Size(),
			Mode:    stat.Mode().String(),
			ModTime: stat.ModTime().Unix(),
			IsDir:   false,
		})
		return files, nil
	}

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, _ := d.Info()
		if info != nil {
			files = append(files, &gpmv1.FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				Mode:    info.Mode().String(),
				ModTime: info.ModTime().Unix(),
				IsDir:   info.IsDir(),
			})
		}

		if d.IsDir() && root != path {
			return filepath.SkipDir
		}

		return nil
	})

	return files, nil
}

func (g *gpm) Pull(ctx context.Context, name string) (<-chan *gpmv1.PullResult, error) {
	ts := func(name string, buf []byte, out chan<- *gpmv1.PullResult) {
		var err error
		var f io.ReadCloser

		defer func() {
			if err != nil {
				out <- &gpmv1.PullResult{Name: name, Error: err.Error()}
			}
		}()

		f, err = os.Open(name)
		if err != nil {
			return
		}
		defer f.Close()

		for {
			nr, er := f.Read(buf)
			if nr > 0 {
				out <- &gpmv1.PullResult{
					Name:     name,
					Chunk:    buf[0:nr],
					Finished: false,
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}

		return
	}

	stat, _ := os.Stat(name)
	if stat == nil {
		return nil, verrs.NotFound(g.Name(), "%s not exist", name)
	}

	buf := make([]byte, 1024*32)
	outs := make(chan *gpmv1.PullResult, 10)

	go func() {
		if !stat.IsDir() {
			ts(name, buf, outs)
		} else {
			_ = filepath.Walk(name, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				ts(path, buf, outs)
				return nil
			})
		}

		outs <- &gpmv1.PullResult{Finished: true}
	}()

	return outs, nil
}

func (g *gpm) Push(ctx context.Context, in <-chan *gpmv1.PushIn) (<-chan *gpmv1.PushResult, error) {
	panic("implement me")
}

func (g *gpm) Exec(ctx context.Context, in *gpmv1.ExecIn) (<-chan *gpmv1.ExecResult, error) {

	cmd := exec.CommandContext(ctx, in.Name, in.Args...)
	execSysProcAttr(cmd, in.Uid, in.Gid)

	out := make(chan *gpmv1.ExecResult, 10)
	ech := make(chan error, 1)
	done := make(chan struct{}, 1)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	buf := make([]byte, 1024*32)
	go func() {
		for {
			select {
			case <-ech:
				stderr.Read(buf)
				out <- &gpmv1.ExecResult{Error: string(buf), Finished: true}
				return
			case <-done:
				stdout.Read(buf)
				out <- &gpmv1.ExecResult{Result: string(buf), Finished: true}
				return
			default:

			}

			stdout.Read(buf)
			out <- &gpmv1.ExecResult{Result: string(buf)}
		}
	}()

	go func() {
		if err = cmd.Wait(); err != nil {
			ech <- err
		} else {
			done <- struct{}{}
		}
	}()

	return out, nil
}

func (g *gpm) Terminal(ctx context.Context, in <-chan *gpmv1.TerminalIn) (<-chan *gpmv1.TerminalResult, error) {
	panic("implement me")
}
