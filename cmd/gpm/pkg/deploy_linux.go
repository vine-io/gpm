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
	"strconv"
	"strings"

	"github.com/vine-io/cli"
	"github.com/vine-io/gpm/pkg/runtime"
	"github.com/vine-io/pkg/release"
)

//go:embed testdata/gpmd
var f embed.FS

const (
	gpm      = "/usr/local/sbin/gpm"
	gpmd     = "/usr/local/sbin/gpmd"
	gpmdConf = "/etc/gpm/gpmd.conf"
)

var gpmdTmp = `#!/bin/bash
# chkconfig: - 30 21
# description: gpmd service.
# Source Function Library
# gpmd Settings

CONF=$(cat /etc/gpm/gpmd.conf)
RETVAL=0
prog="gpmd"

start() {
     echo -n $"Starting $prog: "
     /usr/local/sbin/gpm run $CONF
     RETVAL=$?
     echo
     return $RETVAL
}

stop() {
     echo -n $"Stopping $prog: "
     ps aux | grep "gpmd" | grep -v "grep" | awk -F' ' '{print $2}' |xargs kill -9 
     RETVAL=$?
     echo
     return $RETVAL
 }

restart(){
     stop
     start
}


case "$1" in
 start)
     start
     ;;
 stop)
     stop
     ;;
 restart)
     restart
     ;;
 *)
     echo $"Usage: $0 {start|stop|restart}"
     RETVAL=1
 esac

 exit $RETVAL
`

func deploy(c *cli.Context) error {

	isRun := c.Bool("run")
	args := c.StringSlice("args")

	outE := os.Stdout
	root := "/opt/gpm"
	vv, err := exec.Command("gpm", "-v").CombinedOutput()
	if err == nil {
		fmt.Fprintln(outE, "gpm already exists")
		_ = shutdown(c)
	}
	v := strings.ReplaceAll(string(vv), "\n", "")

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

	// ?????? gpm
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
	_ = dst.Chmod(0o777)

	link, _ := os.Readlink(gpm)
	_ = os.Remove(gpm)
	err = os.Symlink(fname, gpm)
	if err != nil {
		return fmt.Errorf("create gpm link: %v", err)
	}
	if v != "" && v != runtime.GitTag {
		os.Remove(link)
		fmt.Fprintf(outE, "remove old version: %v\n", link)
	}

	// ?????? gpmd
	buf, err := f.ReadFile("testdata/gpmd")
	if err != nil {
		return fmt.Errorf("get gpmd binary: %v", err)
	}
	fname = filepath.Join(root, "bin", "gpmd-"+runtime.GitTag)
	err = ioutil.WriteFile(fname, buf, 0o777)
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
	if v != "" && v != runtime.GitTag {
		os.Remove(link)
		fmt.Fprintf(outE, "remove old version: %v\n", link)
	}

	fmt.Fprintf(outE, "install gpm %s successfully!\n", runtime.GitTag)

	if isRun {
		cmd := exec.Command("gpmd", args...)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("gpmd start: %v", err)
		}

		fmt.Fprintf(outE, "start gpmd successfully!\n")
	}

	_ = os.MkdirAll("/etc/gpm", 0o755)
	text := `--args "--server-address=0.0.0.0:7700" --args "--enable-log"`
	ioutil.WriteFile(gpmdConf, []byte(text), 0o755)

	autoStartFile := "/etc/rc.d/init.d/gpmd"
	_ = ioutil.WriteFile(autoStartFile, []byte(gpmdTmp), 0o777)
	or, e := release.Get()
	if e == nil {
		if (or.Name == "centos" || or.Name == "rhel") && or.Version == "6" {
			exec.Command("chkconfig", "gpmd", "on").CombinedOutput()
		}
	}

	fmt.Fprintf(outE, "install gpm successfully!\n")
	return nil
}

func DeployCmd() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy gpmd and gpm",
		Action: deploy,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "run",
				Usage: "run gpmd after deployed",
			},
			&cli.StringSliceFlag{
				Name:    "args",
				Aliases: []string{"A"},
				Usage:   "the specify args for gpmd",
			},
		},
	}
}

func run(c *cli.Context) error {

	outE := os.Stdout
	args := c.StringSlice("args")
	if len(args) == 0 {
		args = append(args, "--server-address=0.0.0.0:7700", "--enable-log")
	}
	cmd := exec.Command(gpmd, args...)
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("gpmd start: %v", err)
	}

	autoStartFile := "/etc/rc.d/init.d/gpmd"
	stat, _ := os.Stat(autoStartFile)
	if stat != nil {
		text := ""
		for _, arg := range args {
			text += fmt.Sprintf(`--args "%s" `, arg)
		}
		ioutil.WriteFile(gpmdConf, []byte(text), 0o755)
	}

	fmt.Fprintf(outE, "start gpmd successfully!\n")
	return nil
}

func RunCmd() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "run gpmd process",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "args",
				Aliases: []string{"A"},
				Usage:   "the specify args for gpmd",
			},
		},
	}
}

func shutdown(c *cli.Context) error {

	outE := os.Stdout
	b, err := exec.Command("sh", "-c", `ps aux|grep "gpmd" | grep -v "grep" |  awk -F' ' '{print $2}'`).CombinedOutput()
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
	defer p.Wait()

	if err = p.Kill(); err != nil {
		return err
	}

	fmt.Fprintf(outE, "shutdown gpmd successfully!\n")
	return nil
}

func ShutdownCmd() *cli.Command {
	return &cli.Command{
		Name:   "shutdown",
		Usage:  "stop gpmd process",
		Action: shutdown,
	}
}
