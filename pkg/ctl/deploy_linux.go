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

//go:build linux

package ctl

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/wrap"
	"github.com/vine-io/gpm/pkg/server"
	"github.com/vine-io/vine"
	vcmd "github.com/vine-io/vine/lib/cmd"
)

var f embed.FS

const (
	gpm    = "/usr/sbin/gpm"
	gpmCfg = "/etc/gpm.yml"
)

var gpmdTmp = `[Unit]
Description=golang process manager daemon
After=network.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/gpm/
ExecStart=/usr/sbin/gpm run
Restart=on-failure
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
`

var gpmdYaml = `
server:
  address: 0.0.0.0:33700

gpm:
  root: /opt/gpm

logger:
  zap:
    format: json
    filename: /var/log/gpm.log
    max-age: 15
    max-backups: 3
    max-size: 1
`

func deploy(c *cobra.Command, args []string) error {

	outE := os.Stdout
	root := "/opt/gpm"
	vv, err := exec.Command("gpm", "-v").CombinedOutput()
	if err == nil {
		fmt.Fprintln(outE, "gpm already exists")
		_ = shutdown(c, args)
	}
	v := strings.ReplaceAll(strings.TrimPrefix(string(vv), "gpm version "), "\n", "")
	if v != "" {
		fmt.Fprintf(outE, "old gpm: %v\n", v)
	}

	dirs := []string{
		root,
		filepath.Join(root, "bin"),
		filepath.Join(root, "logs"),
		filepath.Join(root, "services"),
		filepath.Join(root, "packages"),
	}

	for _, dir := range dirs {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	// 安装 gpm
	src, err := os.Open(os.Args[0])
	if err != nil {
		return fmt.Errorf("get gpm binary: %v", err)
	}
	defer src.Close()
	fname := filepath.Join(root, "bin", "gpm-"+internal.GitTag)
	dst, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("install gpm: %v", err)
	}
	_ = dst.Chmod(os.ModePerm)

	// 删除旧版本
	link, _ := os.Readlink(gpm)
	if v != "" && v != internal.GitTag {
		os.Remove(link)
		fmt.Fprintf(outE, "remove old version: %v\n", link)
	}
	_ = os.Remove(gpm)

	err = os.Symlink(fname, gpm)
	if err != nil {
		return fmt.Errorf("create gpm link: %v", err)
	}

	autoStartFile := "/usr/lib/systemd/system/gpm.service"
	_ = os.WriteFile(autoStartFile, []byte(gpmdTmp), os.ModePerm)
	exec.Command("systemctl", "enable", "service").CombinedOutput()
	exec.Command("systemctl", "daemon-reload").CombinedOutput()

	_ = os.WriteFile(gpmCfg, []byte(gpmdYaml), os.ModePerm)

	fmt.Fprintf(outE, "install gpm successfully!\n")
	return nil
}

func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "deploy gpmd and gpm",
		GroupID: "gpm",
		RunE:    deploy,
	}

	return cmd
}

func initRun(root *cobra.Command) error {

	s := vine.NewService(
		vine.Cmd(vcmd.NewCmd(vcmd.NewApp(root))),
		vine.Name(internal.GpmName),
		vine.ID(internal.GpmId),
		vine.Address(":33700"),
		vine.Version(internal.GetVersion()),
		vine.Metadata(map[string]string{
			"namespace": internal.Namespace,
		}),
		vine.WrapHandler(wrap.NewLoggerWrapper()),
	)

	app, err := server.New(s)
	if err != nil {
		return err
	}

	action := vine.Action(func(cmd *cobra.Command, args []string) error {
		if err = server.Action(cmd, args); err != nil {
			return err
		}

		if err = app.Init(); err != nil {
			return err
		}

		return app.Run()
	})

	return s.Init(action)
}

func RunCmd() (*cobra.Command, error) {
	runCmd := &cobra.Command{
		Use:     "run",
		Short:   "run gpmd process",
		GroupID: "gpm",
	}

	if err := initRun(runCmd); err != nil {
		return nil, err
	}

	return runCmd, nil
}

func shutdown(c *cobra.Command, args []string) error {

	outE := os.Stdout
	b, err := exec.Command("bash", "-c", `ps aux|grep "gpm run" | grep -v "grep" | head -1 |  awk -F' ' '{print $2}'`).CombinedOutput()
	if err != nil {
		return fmt.Errorf("get gpmd failed: %v", err)
	}
	pid := strings.ReplaceAll(string(b), "\n", "")

	fmt.Fprintf(outE, "gpmd pid=%s\n", pid)
	Pid, _ := strconv.ParseInt(pid, 10, 64)

	p, err := os.FindProcess(int(Pid))
	if err != nil {
		return fmt.Errorf("get gpmd failed: %v", err)
	}

	if err = p.Kill(); err != nil {
		return err
	}

	_ = p.Release()

	for {
		state, _ := p.Wait()
		if state == nil || state.Exited() {
			break
		}
		time.Sleep(time.Second)
	}

	fmt.Fprintf(outE, "shutdown gpmd successfully!\n")
	return nil
}

func ShutdownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shutdown",
		Short:   "stop gpmd process",
		GroupID: "gpm",
		RunE:    shutdown,
	}

	return cmd
}
