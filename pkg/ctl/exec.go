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

package ctl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vine-io/cli"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal/client"
	verrs "github.com/vine-io/vine/lib/errors"
)

func execBash(c *cli.Context) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	in := &gpmv1.ExecIn{}
	in.Shell = c.String("shell")
	in.Dir = c.String("dir")
	env := c.StringSlice("env")
	in.User = c.String("user")
	in.Group = c.String("group")
	if err := in.Validate(); err != nil {
		return err
	}

	if len(env) > 0 {
		in.Env = map[string]string{}
	}
	for _, item := range env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			in.Env[parts[0]] = parts[1]
		}
	}

	result, err := cc.Exec(ctx, in, opts...)
	if err != nil {
		return fmt.Errorf("%v", verrs.FromErr(err).Detail)
	}

	fmt.Fprintln(outE, string(result.Result))

	return nil
}

func ExecBashCmd() *cli.Command {
	return &cli.Command{
		Name:     "exec",
		Usage:    "execute command",
		Category: "bash",
		Action:   execBash,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "shell",
				Aliases: []string{"S"},
				Usage:   "specify the command for exec",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "specify the directory path for exec",
			},
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"E"},
				Usage:   "specify the env for exec",
			},
			&cli.StringFlag{
				Name:  "user",
				Usage: "specify the user for exec",
			},
			&cli.StringFlag{
				Name:  "group",
				Usage: "specify the group for exec",
			},
		},
	}
}
