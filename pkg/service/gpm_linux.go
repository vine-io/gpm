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

package service

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
)

func fillService(service *gpmv1.Service) error {
	fn := func(attr *gpmv1.SysProcAttr) {
		u, err := user.Current()
		if err == nil {
			User, _ := user.LookupId(u.Uid)
			Group, _ := user.LookupGroupId(u.Gid)
			attr.User = User.Name
			attr.Group = Group.Name
			uid, _ := strconv.ParseInt(u.Uid, 10, 64)
			gid, _ := strconv.ParseInt(Group.Gid, 10, 64)
			attr.Uid = int32(uid)
			attr.Gid = int32(gid)
		}
	}

	attr := service.SysProcAttr
	if attr == nil {
		attr = &gpmv1.SysProcAttr{}
		fn(attr)
	} else {
		if attr.User == "" || attr.Group == "" {
			fn(attr)
		}
	}

	if attr.Uid == 0 {
		u, err := user.Lookup(attr.User)
		if err != nil {
			return err
		}
		uid, _ := strconv.ParseInt(u.Uid, 10, 64)
		attr.Uid = int32(uid)
	}
	if attr.Gid == 0 {
		group, err := user.LookupGroup(attr.Group)
		if err != nil {
			return err
		}

		gid, _ := strconv.ParseInt(group.Gid, 10, 64)
		attr.Gid = int32(gid)
	}

	return nil
}

func injectSysProcAttr(cmd *exec.Cmd, attr *gpmv1.SysProcAttr) {
	sysAttr := &syscall.SysProcAttr{
		Setpgid: true,
	}

	if attr != nil {
		sysAttr.Credential = &syscall.Credential{
			Uid: uint32(attr.Uid),
			Gid: uint32(attr.Gid),
		}
		if attr.Chroot != "" {
			sysAttr.Chroot = attr.Chroot
		}
	}

	cmd.SysProcAttr = sysAttr
}

func execSysProcAttr(cmd *exec.Cmd, uid, gid int32) {
	sysAttr := &syscall.SysProcAttr{
		Setpgid: true,
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	cmd.SysProcAttr = sysAttr
}

func startTerminal(uid, gid int32, env map[string]string) *exec.Cmd {
	cmd := exec.Command("/bin/bash")

	sysAttr := &syscall.SysProcAttr{
		Setpgid: true,
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	cmd.SysProcAttr = sysAttr
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	home, err := os.UserHomeDir()
	if err == nil {
		cmd.Dir = home
	}

	return cmd
}
