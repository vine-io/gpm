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

// +build linux

package pkg

import (
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/lack-io/cli"
)

//go:embed testdata/gpmd
var f embed.FS

func deploy(c *cli.Context) error {

	outE := os.Stdout
	root := "/opt/gpm"
	_, err := exec.Command("gpm", "-v").CombinedOutput()
	if err == nil {
		return fmt.Errorf("gpm installed")
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
	fname := filepath.Join(root, "bin", "gpm-"+runtime.GitTag)
	dst, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("install gpm: %v", err)
	}
	_ = dst.Chmod(0777)
	err = os.Symlink(fname, "/usr/local/sbin/gpm")
	if err != nil {
		return fmt.Errorf("create gpm link: %v", err)
	}

	// 安装 gpmd
	buf, err := f.ReadFile("testdata/gpmd")
	if err != nil {
		return fmt.Errorf("get gpmd binary: %v", err)
	}
	fname = filepath.Join(root, "bin", "gpmd-"+runtime.GitTag)
	err = ioutil.WriteFile(fname, buf, 0777)
	if err != nil {
		return fmt.Errorf("install gpmd: %v", err)
	}
	err = os.Symlink(fname, "/usr/local/sbin/gpmd")
	if err != nil {
		return fmt.Errorf("create gpmd link: %v", err)
	}

	fmt.Fprintf(outE, "install gpm successfully!\n")
	return nil
}

func DeployCmd() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy gpmd and gpm",
		Action: deploy,
	}
}

func update(c *cli.Context) error {

	fmt.Println(os.Args[0])
	return nil
}

func UpdateCmd() *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "update gpmd and gpm",
		Action: update,
	}
}

