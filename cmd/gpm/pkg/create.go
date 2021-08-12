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

package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vine-io/gpm/pkg/runtime/client"
	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
	"github.com/vine-io/cli"
)

func createService(c *cli.Context) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout
	spec := &gpmv1.ServiceSpec{
		Env:         map[string]string{},
		SysProcAttr: &gpmv1.SysProcAttr{},
		Log:         &gpmv1.ProcLog{},
		AutoRestart: 1,
	}

	spec.Name = c.String("name")
	spec.Bin = c.String("bin")
	spec.Args = c.StringSlice("args")
	spec.Dir = c.String("dir")
	env := c.StringSlice("env")
	spec.SysProcAttr.User = c.String("user")
	spec.SysProcAttr.Group = c.String("group")
	spec.Log.Expire = int32(c.Int("log-expire"))
	spec.Log.MaxSize = c.Int64("log-max-size")
	spec.Version = c.String("version")
	autoRestart := c.Bool("auto-restart")
	if err := spec.Validate(); err != nil {
		return err
	}
	for _, item := range env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			spec.Env[parts[0]] = parts[1]
		}
	}
	if autoRestart {
		spec.AutoRestart = 1
	} else {
		spec.AutoRestart = 0
	}

	s, err := cc.CreateService(ctx, spec, opts...)
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "service '%s' %s\n", s.Name, s.Status)
	return nil
}

func CreateServiceCmd() *cli.Command {
	return &cli.Command{
		Name:     "create",
		Usage:    "create a service",
		Category: "service",
		Action:   createService,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"N"},
				Usage:   "specify the name for service",
			},
			&cli.StringFlag{
				Name:    "bin",
				Aliases: []string{"B"},
				Usage:   "specify the bin for service",
			},
			&cli.StringSliceFlag{
				Name:    "args",
				Aliases: []string{"A"},
				Usage:   "specify the args for service",
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"D"},
				Usage:   "specify the root directory for service",
			},
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"E"},
				Usage:   "specify the env for service",
			},
			&cli.StringFlag{
				Name:  "user",
				Usage: "specify the user for service",
			},
			&cli.StringFlag{
				Name:  "group",
				Usage: "specify the group for service",
			},
			&cli.IntFlag{
				Name:  "log-expire",
				Usage: "specify the expire for service log",
				Value: 15,
			},
			&cli.Int64Flag{
				Name:  "log-max-size",
				Usage: "specify the max size for service log",
				Value: 1024 * 1024 * 10,
			},
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "specify the version for service",
			},
			&cli.BoolFlag{
				Name:  "auto-restart",
				Usage: "Whether auto restart service when it crashing",
				Value: true,
			},
		},
	}
}
