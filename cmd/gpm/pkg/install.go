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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	pbr "github.com/schollz/progressbar/v3"
	"github.com/vine-io/cli"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/runtime/client"
	"google.golang.org/grpc/status"
)

func installService(c *cli.Context) error {

	opts := getCallOptions(c)
	spec := &gpmv1.ServiceSpec{
		Env:         map[string]string{},
		SysProcAttr: &gpmv1.SysProcAttr{},
		Log:         &gpmv1.ProcLog{},
		AutoRestart: 1,
	}

	pack := c.String("package")
	if len(pack) == 0 {
		return fmt.Errorf("missing package")
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

	cc := client.New()
	ctx := context.Background()
	ech := make(chan error, 1)
	done := make(chan struct{}, 1)
	buf := make([]byte, 1024*32)
	outE := os.Stdout

	file, err := os.Open(pack)
	if err != nil {
		return err
	}

	s, err := cc.InstallService(ctx, spec, opts...)
	if err != nil {
		return err
	}

	pb := pbr.NewOptions(0,
		pbr.OptionSetWriter(outE),
		pbr.OptionSetDescription(fmt.Sprintf("upload [%s]", pack)),
		pbr.OptionShowBytes(true),
		pbr.OptionEnableColorCodes(true),
		pbr.OptionOnCompletion(func() {
			fmt.Fprintf(outE, "\n")
		}),
	)

	go func() {
		p := &gpmv1.Package{Package: filepath.Base(pack)}
		stat, _ := file.Stat()
		if stat != nil {
			p.Total = stat.Size()
		}
		pb.ChangeMax64(p.Total)
		for {
			n, e := file.Read(buf)
			if e != nil && e != io.EOF {
				ech <- e
				return
			}

			if e == io.EOF {
				p.IsOk = true
			}

			_ = pb.Add(n)
			p.Length = int64(n)
			p.Chunk = buf[0:n]
			e1 := s.Send(p)
			if e1 != nil {
				return
			}

			if e == io.EOF {
				break
			}
		}
	}()

	go func() {
		for {
			b, err := s.Recv()
			if err != nil {
				ech <- errors.New(status.Convert(err).Message())
				return
			}
			if len(b.Error) != 0 {
				ech <- errors.New(b.Error)
				return
			}
			if b.IsOk {
				done <- struct{}{}
				return
			}
		}
	}()

	select {
	case e := <-ech:
		_ = pb.Clear()
		return e
	case <-done:
	}

	fmt.Fprintf(outE, "install service %s successfully\n", spec.Name)
	return nil
}

func InstallServiceCmd() *cli.Command {
	return &cli.Command{
		Name:     "install",
		Usage:    "install a service",
		Category: "service",
		Action:   installService,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"P"},
				Usage:   "specify the package for service",
			},
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
