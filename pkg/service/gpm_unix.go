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
	"os/user"
	"strconv"

	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
)

func fillService(service *gpmv1.Service) error {
	if service.User == "" {
		service.User = "root"
	}
	if service.Group == "" {
		service.Group = "root"
	}

	u, err := user.Lookup(service.User)
	if err != nil {
		return err
	}

	group, err := user.LookupGroup(service.Group)
	if err != nil {
		return err
	}

	uid, _ := strconv.ParseInt(u.Uid, 10, 64)
	gid, _ := strconv.ParseInt(group.Gid, 10, 64)
	service.Uid = int32(uid)
	service.Gid = int32(gid)

	return nil
}
