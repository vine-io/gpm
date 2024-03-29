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

//go:build windows

package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func fillService(service *gpmv1.Service) error {
	return nil
}

func injectSysProcAttr(cmd *exec.Cmd, attr *gpmv1.SysProcAttr) {
	sysAttr := &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.SysProcAttr = sysAttr
}

func execSysProcAttr(cmd *exec.Cmd, in *gpmv1.ExecIn) {

	sysAttr := &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.Env = os.Environ()
	for k, v := range in.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.SysProcAttr = sysAttr
	cmd.Dir = in.Dir
}

func adminCmd(cmd *exec.Cmd) {
	sysAttr := &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.Env = os.Environ()
	cmd.SysProcAttr = sysAttr
}

func startExec(ctx context.Context, in *gpmv1.ExecIn) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "cmd", "/C", in.Shell)

	sysAttr := &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.SysProcAttr = sysAttr
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range in.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	home, err := os.UserHomeDir()
	if err == nil {
		cmd.Dir = home
	}

	return cmd
}

func startTerminal(ctx context.Context, in *gpmv1.TerminalIn) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "powershell.exe")

	sysAttr := &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.SysProcAttr = sysAttr
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range in.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	home, err := os.UserHomeDir()
	if err == nil {
		cmd.Dir = home
	}

	return cmd
}

func beauty(b []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(bytes.TrimSuffix(b, []byte("\r\n"))), simplifiedchinese.GBK.NewDecoder())
	d, _ := io.ReadAll(reader)
	return d
}
