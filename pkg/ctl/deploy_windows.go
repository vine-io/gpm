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
// +build windows

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
)

//go:embed testdata/nssm.exe
//go:embed testdata/gpmd.exe
var f embed.FS

const (
	root = "C:\\opt\\gpm"
	gpm  = "c:\\opt\\gpm\\bin\\gpm.exe"
	gpmd = "c:\\opt\\gpm\\bin\\gpmd.exe"
)

func deploy(c *cobra.Command, args []string) error {

	isRun, _ := c.Flags().GetBool("run")
	dArgs, _ := c.Flags().GetStringSlice("args")

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
		_ = os.MkdirAll(dir, 0o777)
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
	_ = dst.Chmod(0o777)

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

	// 安装 gpmd
	buf, err := f.ReadFile("testdata/gpmd.exe")
	if err != nil {
		return fmt.Errorf("get gpmd binary: %v", err)
	}
	fname = filepath.Join(root, "bin", "gpmd-"+internal.GitTag+".exe")
	err = os.WriteFile(fname, buf, 0o777)
	if err != nil {
		return fmt.Errorf("install gpmd: %v", err)
	}
	_ = os.Chmod(fname, 0o777)

	link, _ = os.Readlink(gpmd)
	_ = os.Remove(gpmd)
	err = os.Symlink(fname, gpmd)
	if err != nil {
		return fmt.Errorf("create gpmd link: %v", err)
	}
	if filepath.Base(link) != "gpmd"+"-"+internal.GitTag+".exe" {
		os.Remove(link)
		fmt.Fprintf(outE, "remove old version: %v\n", link)
	}

	// 安装 nssm.exe
	stat, _ := os.Stat(root + "\\bin\\nssm.exe")
	if stat == nil {
		buf, err = f.ReadFile("testdata/nssm.exe")
		if err != nil {
			return fmt.Errorf("get nssm binary: %v", err)
		}
		fname = filepath.Join(root, "bin", "nssm.exe")
		err = ioutil.WriteFile(fname, buf, 0o777)
		if err != nil {
			return fmt.Errorf("install nssm: %v", err)
		}
		_ = os.Chmod(fname, 0o777)
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
	nssmShell := fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe install gpmd %s\\bin\\gpmd", root, root)
	if len(dArgs) > 0 {
		nssmShell += " " + strings.Join(dArgs, " ")
	}
	bb.WriteString(fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe remove gpmd confirm\r\n", root))
	bb.WriteString(nssmShell)
	startBat := filepath.Join(os.TempDir(), "start.bat")
	defer os.Remove(startBat)
	_ = ioutil.WriteFile(startBat, bb.Bytes(), 0o777)
	_, err = exec.Command("cmd", "/C", startBat).CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute %s: %v", startBat, err)
	}

	if isRun {
		cmd := exec.Command(root+"\\bin\\"+"nssm.exe", "start", "gpmd")
		if _, err = cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("gpmd start: %v", err)
		}

		fmt.Fprintf(outE, "start gpmd successfully!\n")
	}

	fmt.Fprintf(outE, "install gpm successfully!\n")
	return nil
}

func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy gpmd and gpm",
		RunE:  deploy,
	}

	cmd.PersistentFlags().Bool("run", false, "run gpmd after deployed")
	cmd.PersistentFlags().StringP("args", "A", "", "the specify args for gpmd")

	return cmd
}

func run(c *cobra.Command, args []string) error {

	outE := os.Stdout
	args, _ = c.Flags().GetStringSlice("args")

	bb := bytes.NewBuffer([]byte("@echo off\r\n\r\n"))
	// TODO: 设置 nssm
	nssmShell := fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe install gpmd %s\\bin\\gpmd", root, root)
	if len(args) > 0 {
		nssmShell += " " + strings.Join(args, " ")
	}
	bb.WriteString(fmt.Sprintf("\r\n\r\n%s\\bin\\nssm.exe remove gpmd confirm\r\n", root))
	bb.WriteString(nssmShell)
	startBat := filepath.Join(root, "start.bat")
	defer os.Remove(startBat)
	_ = ioutil.WriteFile(startBat, bb.Bytes(), 0o777)
	_, err := exec.Command("cmd", "/C", startBat).CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute %s: %v", startBat, err)
	}

	cmd := exec.Command(root+"\\bin\\"+"nssm.exe", "start", "gpmd")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("gpmd start: %v", err)
	}

	fmt.Fprintf(outE, "start gpmd successfully!\n")
	return nil
}

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run gpmd process",
		RunE:  run,
	}

	cmd.PersistentFlags().StringSliceP("args", "A", []string{}, "the specify args for gpmd")

	return cmd
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
		Use:   "shutdown",
		Short: "stop gpmd process",
		RunE:  shutdown,
	}

	return cmd
}
