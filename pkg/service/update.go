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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	gruntime "runtime"
	"strings"
	"time"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	log "github.com/lack-io/vine/lib/logger"
)

func (g *gpm) UpdateSelf(ctx context.Context, version string, in <-chan *gpmv1.UpdateIn) (<-chan *gpmv1.UpdateResult, error) {

	log.Infof("update gpm and gpmd to %s", version)
	outs := make(chan *gpmv1.UpdateResult, 10)
	go func() {
		//if runtime.GitTag == version {
		//	outs <- &gpmv1.UpdateResult{Error: "consistent version: " + version}
		//	return
		//}
		// 接收 gpm 文件
		dst := filepath.Join(os.TempDir(), fmt.Sprintf("gpm-tmp"))
		if gruntime.GOOS == "windows" {
			dst += ".exe"
		}
		f, err := os.Create(dst)
		if err != nil {
			outs <- &gpmv1.UpdateResult{Error: err.Error()}
			return
		}
		defer f.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case b, ok := <-in:
				if !ok {
					return
				}

				if b.Length > 0 {
					_, err = f.Write(b.Chunk)
					if err != nil {
						outs <- &gpmv1.UpdateResult{Error: err.Error()}
						return
					}
				}
				if b.IsOk {
					outs <- &gpmv1.UpdateResult{IsOk: true}
				}

				if b.DeploySignal {
					goto EXIT
				}
			}
		}
	EXIT:
		f.Chmod(0777)
		_ = f.Close()

		args := []string{"deploy", "--run", "--args", fmt.Sprintf(`"--server-address=%s"`, g.Cfg.Address)}
		if g.Cfg.EnableLog {
			args = append(args, fmt.Sprintf(`--args "--enable-log"`))
		}
		shell := fmt.Sprintf(`%s %s`, dst, strings.Join(args, " "))
		log.Infof("starting upgrade gpmd and gpm, %s", shell)
		script := filepath.Join(os.TempDir(), "start.sh")

		_ = ioutil.WriteFile(script, []byte(fmt.Sprintf(`#! /bin/bash
nohup %s &`, shell)), 0777)
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

	return outs, nil
}
