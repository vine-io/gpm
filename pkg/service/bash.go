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
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	log "github.com/lack-io/vine/lib/logger"
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

		if root == path {
			return nil
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

func (g *gpm) Pull(ctx context.Context, name string, isDir bool, sender Sender) error {
	ts := func(name string, total int64, buf []byte, stream Sender) error {
		var err error
		var f *os.File
		out := &gpmv1.PullResult{}

		defer func() {
			if err != nil {
				out.Name, out.Error = name, err.Error()
				_ = stream.Send(out)
			}
		}()

		f, err = os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		for {
			nr, er := f.Read(buf)
			if nr > 0 {
				out = &gpmv1.PullResult{
					Name:     name,
					Total:    total,
					Chunk:    buf[0:nr],
					Length:   int64(nr),
					Finished: false,
				}
				stream.Send(out)
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}

		return err
	}

	stat, _ := os.Stat(name)
	if stat == nil {
		return verrs.NotFound(g.Name(), "%s not exist", name)
	}
	if stat.IsDir() && !isDir {
		return verrs.BadRequest(g.Name(), "%s is a directory", name)
	}
	if !stat.IsDir() && isDir {
		return verrs.BadRequest(g.Name(), "%s is not a directory", name)
	}

	buf := make([]byte, 1024*32)

	var err error
	if !stat.IsDir() {
		err = ts(name, stat.Size(), buf, sender)
	} else {
		err = filepath.Walk(name, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			return ts(path, info.Size(), buf, sender)
		})
	}
	if err != nil {
		return err
	}

	return sender.Send(&gpmv1.PullResult{Finished: true})
}

func (g *gpm) Push(ctx context.Context, stream Stream) error {
	var file *os.File
	out := &gpmv1.PushResult{}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Infof("push stream done!")
			return nil
		default:
		}

		data, err := stream.Recv()
		if err != nil && err != io.EOF {
			out.Error = err.Error()
			_ = stream.Send(out)
			return err
		}
		b := data.(*gpmv1.PushIn)

		if file == nil {
			if err = b.Validate(); err != nil {
				out.Error = verrs.BadRequest(g.Name(), err.Error()).Error()
				_ = stream.Send(out)
				return err
			}

			dst := b.Dst
			stat, _ := os.Stat(b.Dst)
			if stat != nil && stat.IsDir() {
				dst = filepath.Join(dst, b.Name)
			}
			dir := filepath.Dir(dst)
			if stat == nil {
				err = os.MkdirAll(dir, 0o777)
				if err != nil {
					out.Error = verrs.InternalServerError(g.Name(), err.Error()).Error()
					_ = stream.Send(out)
					return err
				}
			}

			log.Infof("save file: %s -> %s", b.Name, dst)
			file, err = os.Create(dst)
			if err != nil {
				out.Error = verrs.InternalServerError(g.Name(), err.Error()).Error()
				_ = stream.Send(out)
				return err
			}
		}

		_, err = file.Write(b.Chunk[0:b.Length])
		if err != nil {
			out.Error = err.Error()
			_ = stream.Send(out)
			return err
		}
		if b.IsOk {
			out.IsOk = true
			_ = stream.Send(out)
			return nil
		}
	}
}

func (g *gpm) Exec(ctx context.Context, in *gpmv1.ExecIn) (*gpmv1.ExecResult, error) {

	cmd := exec.CommandContext(ctx, in.Name, in.Args...)
	execSysProcAttr(cmd, in)

	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, verrs.InternalServerError(g.Name(), "exec %w: %v", err, beauty(b))
	}

	out := &gpmv1.ExecResult{
		Result: beauty(b),
	}

	return out, nil
}

type wr struct {
	stream Stream
}

func (wr *wr) Write(data []byte) (int, error) {
	err := wr.stream.Send(&gpmv1.TerminalResult{Stdout: data})
	return len(data), err
}

func (wr *wr) Read(data []byte) (int, error) {
	buf, err := wr.stream.Recv()
	if err != nil {
		return 0, err
	}
	b := buf.(*gpmv1.TerminalIn)
	return bytes.NewBuffer(data).WriteString(b.Command)
}

func (g *gpm) Terminal(ctx context.Context, stream Stream) error {
	data, err := stream.Recv()
	if err != nil {
		return err
	}
	b := data.(*gpmv1.TerminalIn)

	if err = b.Validate(); err != nil {
		return verrs.BadGateway(g.Name(), err.Error())
	}
	tid := uuid.New().String()

	into := &wr{stream: stream}
	cmd := startTerminal(b)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}

	go io.Copy(stdin, into)
	go io.Copy(into, stdout)
	go io.Copy(into, stderr)

	if err = cmd.Start(); err != nil {
		return verrs.InternalServerError(g.Name(), err.Error())
	}

	_, err = cmd.Process.Wait()
	log.Infof("terminal %s done!", tid)

	return err
}
