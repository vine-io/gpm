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

// +build !windows

package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
	log "github.com/vine-io/vine/lib/logger"
	verrs "github.com/vine-io/vine/proto/apis/errors"
)

func (g *gpm) UpdateSelf(ctx context.Context, stream IOStream) error {

	var (
		file *os.File
		err  error
		dst  string
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
		b := data.(*gpmv1.UpdateIn)

		if file == nil {
			log.Infof("update gpm and gpmd to %s", b.Version)
			if err = b.Validate(); err != nil {
				return verrs.BadRequest(g.Name(), err.Error())
			}

			// 接收 gpm 文件
			dst = filepath.Join(os.TempDir(), fmt.Sprintf("gpm-tmp"))
			file, err = os.Create(dst)
			if err != nil {
				//outs <- &gpmv1.UpdateResult{Error: err.Error()}
				return err
			}
		}

		if b.Length > 0 {
			_, err = file.Write(b.Chunk)
			if err != nil {
				//outs <- &gpmv1.UpdateResult{Error: err.Error()}
				return err
			}
		}
		if b.IsOk {
			stream.Send(&gpmv1.UpdateResult{IsOk: true})
		}

		if b.DeploySignal {
			goto EXIT
		}
	}

EXIT:
	file.Chmod(0o777)
	_ = file.Close()

	go func() {
		args := []string{"deploy", "--run", "--args", fmt.Sprintf(`"--server-address=%s"`, g.Cfg.Get("server", "address").String(""))}
		if g.Cfg.Get("enable", "log").Bool(false) {
			args = append(args, fmt.Sprintf(`--args "--enable-log"`))
		}
		shell := fmt.Sprintf(`%s %s`, dst, strings.Join(args, " "))
		log.Infof("starting upgrade gpmd and gpm, %s", shell)
		script := filepath.Join(os.TempDir(), "start.sh")

		_ = ioutil.WriteFile(script, []byte(fmt.Sprintf(`#! /bin/bash
nohup %s &`, shell)), 0o777)
		cmd := exec.Command("/bin/bash", "-c", shell)
		adminCmd(cmd)
		err = cmd.Start()
		if err != nil {
			log.Errorf("run gpmd: %v", err)
		}
		var pr *os.Process
		for {
			pr = cmd.Process
			if pr != nil {
				log.Infof("start gpm at %d", pr.Pid)
				break
			}
			time.Sleep(time.Second * 1)
		}

		log.Infof("gpmd shutdown, waiting for reboot")
		os.Exit(0)
	}()

	return nil
}
