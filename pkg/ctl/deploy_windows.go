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

package ctl

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/client"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/wrap"
	"github.com/vine-io/gpm/pkg/server"
	"github.com/vine-io/vine"
	vcmd "github.com/vine-io/vine/lib/cmd"
)

//go:embed testdata/nssm.exe
var f embed.FS

const (
	root = "C:\\opt\\gpm"
	gpm  = "c:\\opt\\gpm\\bin\\gpm.exe"
)

var gpmdYaml = `
server:
  address: 0.0.0.0:33700

gpm:
  root: /opt/gpm
`

func deploy(c *cobra.Command, args []string) error {

	outE := os.Stdout
	_, err := exec.Command(root+"\\bin\\"+"gpm", "-v").CombinedOutput()
	if err == nil {
		fmt.Fprintln(outE, "gpm already exists")
		_ = shutdown(c, args)
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
	fname := filepath.Join(root, "bin", "gpm-"+internal.GitTag+".exe")
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

	link, _ := os.Readlink(gpm)
	_ = os.Remove(gpm)
	err = os.Symlink(fname, gpm)
	if err != nil {
		return fmt.Errorf("create gpm link: %v", err)
	}
	if filepath.Base(link) != "gpm"+"-"+internal.GitTag+".exe" {
		os.Remove(link)
		fmt.Fprintf(outE, "remove old version: %v\n", link)
	}

	// 安装 nssm.exe
	stat, _ := os.Stat(root + "\\bin\\nssm.exe")
	if stat == nil {
		buf, err := f.ReadFile("testdata/nssm.exe")
		if err != nil {
			return fmt.Errorf("get nssm binary: %v", err)
		}
		fname = filepath.Join(root, "bin", "nssm.exe")
		err = ioutil.WriteFile(fname, buf, os.ModePerm)
		if err != nil {
			return fmt.Errorf("install nssm: %v", err)
		}
		_ = os.Chmod(fname, os.ModePerm)
	}

	fmt.Fprintf(outE, "install gpm %s successfully!\n", internal.GitTag)

	bb := bytes.NewBuffer([]byte("@echo off\r\n\r\n"))
	path := os.Getenv("Path")
	if path != "" {
		if !strings.Contains(path, root+"\\bin") {
			bb.WriteString(fmt.Sprintf(`setx "Path" "%s\bin;%s" /m`, root, path))
		}
	}

	// TODO: 设置 nssm
	dArgs := []string{"run"}
	nssmShell := fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe install gpmd %s\\bin\\gpm", root, root)
	if len(dArgs) > 0 {
		nssmShell += " " + strings.Join(dArgs, " ")
	}
	bb.WriteString(fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe remove gpmd confirm\r\n", root))
	bb.WriteString(nssmShell)
	startBat := filepath.Join(os.TempDir(), "start.bat")
	defer os.Remove(startBat)
	_ = os.WriteFile(startBat, bb.Bytes(), os.ModePerm)
	_, err = exec.Command("cmd", "/C", startBat).CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute %s: %v", startBat, err)
	}

	fmt.Fprintf(outE, "install gpm successfully!\n")
	return nil
}

func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "deploy",
		Short:        "deploy gpmd and gpm",
		GroupID:      "gpm",
		SilenceUsage: true,
		RunE:         deploy,
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
		Use:          "run",
		Short:        "run gpmd process",
		SilenceUsage: true,
		GroupID:      "gpm",
	}

	if err := initRun(runCmd); err != nil {
		return nil, err
	}

	return runCmd, nil
}

func shutdown(c *cobra.Command, args []string) error {

	outE := os.Stdout
	n := 0
	for {
		_, _ = exec.Command(root+"\\bin\\"+"nssm.exe", "stop", "gpmd").CombinedOutput()
		ctx := context.Background()
		cc := client.New()
		if err := cc.Healthz(ctx, getCallOptions(c)...); err != nil {
			break
		}
		if n > 3 {
			return fmt.Errorf("shutdown gpmd timeout")
		}
		time.Sleep(time.Second * 1)
		n += 1
	}

	fmt.Fprintf(outE, "shutdown gpmd successfully!\n")
	return nil
}

func ShutdownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "shutdown",
		Short:        "stop gpmd process",
		GroupID:      "gpm",
		SilenceUsage: true,
		RunE:         shutdown,
	}

	return cmd
}
