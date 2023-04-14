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
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	vserver "github.com/vine-io/vine/core/server"
	verrs "github.com/vine-io/vine/lib/errors"
	log "github.com/vine-io/vine/lib/logger"
)

func NewSFtpService(ctx context.Context, server vserver.Server) (GenerateFTP, error) {

	ftp := &sftp{
		ctx:    ctx,
		server: server,
	}

	return ftp, nil
}

type sftp struct {
	ctx context.Context

	server vserver.Server
}

func (s *sftp) Name() string {
	return s.server.Options().Name
}

func (s *sftp) List(ctx context.Context, root string) ([]*gpmv1.FileInfo, error) {
	stat, _ := os.Stat(root)
	if stat == nil {
		return nil, verrs.NotFound(s.Name(), "invalid path '%s'", root)
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

func (s *sftp) Pull(ctx context.Context, name string, isDir bool, sender IOWriter) error {
	ts := func(name string, total int64, buf []byte, stream IOWriter) error {
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
				_ = stream.Send(out)
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
		return verrs.NotFound(s.Name(), "%s not exist", name)
	}
	if stat.IsDir() && !isDir {
		return verrs.BadRequest(s.Name(), "%s is a directory", name)
	}
	if !stat.IsDir() && isDir {
		return verrs.BadRequest(s.Name(), "%s is not a directory", name)
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

func (s *sftp) Push(ctx context.Context, stream IOReader) error {
	var file *os.File

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	var err error
	for {
		select {
		case <-ctx.Done():
			log.Infof("push stream done!")
			return nil
		default:
		}

		data, e := stream.Recv()
		if e != nil && e != io.EOF {
			return fmt.Errorf("receive data: %v", e)
		}
		b := data.(*gpmv1.PushIn)

		if file == nil {
			if err = b.Validate(); err != nil {
				return err
			}

			dst := b.Dst
			stat, _ := os.Stat(b.Dst)
			if stat != nil && stat.IsDir() {
				dst = filepath.Join(dst, b.Name)
			}
			dir := filepath.Dir(dst)
			if stat == nil {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					return err
				}
			}

			log.Infof("save file: %s -> %s", b.Name, dst)
			file, err = createFile(dst)
			if err != nil {
				return err
			}
		}

		if b.Length > 0 {
			_, err = file.Write(b.Chunk[0:b.Length])
			if err != nil {
				return err
			}
		}

		if b.IsOk {
			break
		}
	}

	return stream.Close()
}

func (s *sftp) Exec(ctx context.Context, in *gpmv1.ExecIn) (*gpmv1.ExecResult, error) {

	cmd := startExec(ctx, in)
	execSysProcAttr(cmd, in)

	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, verrs.InternalServerError(s.Name(), "exec %v: %v", err, string(beauty(b)))
	}

	out := &gpmv1.ExecResult{
		Result: beauty(b),
	}

	return out, nil
}

type wr struct {
	stream IOStream
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

func (s *sftp) Terminal(ctx context.Context, stream IOStream) error {
	data, err := stream.Recv()
	if err != nil {
		return err
	}
	b := data.(*gpmv1.TerminalIn)

	if err = b.Validate(); err != nil {
		return verrs.BadRequest(s.Name(), err.Error())
	}
	tid := uuid.New().String()

	into := &wr{stream: stream}
	cmd := startTerminal(ctx, b)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	go io.Copy(stdin, into)
	go io.Copy(into, stdout)
	go io.Copy(into, stderr)

	if err = cmd.Start(); err != nil {
		return verrs.InternalServerError(s.Name(), err.Error())
	}

	_, err = cmd.Process.Wait()
	log.Infof("terminal %s done!", tid)

	return err
}

func (s *sftp) String() string {
	return "sftp"
}
