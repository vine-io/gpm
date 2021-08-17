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

	pbr "github.com/schollz/progressbar/v3"
	"github.com/vine-io/cli"
	"github.com/vine-io/gpm/pkg/runtime/client"
	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
	"google.golang.org/grpc/status"
)

func upgradeService(c *cli.Context) error {

	pack := c.String("package")
	if len(pack) == 0 {
		return fmt.Errorf("missing package")
	}

	name := c.String("name")
	v := c.String("version")
	if name == "" {
		return fmt.Errorf("missing name")
	}
	if v == "" {
		return fmt.Errorf("missing version")
	}

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	ech := make(chan error, 1)
	done := make(chan struct{}, 1)
	buf := make([]byte, 1024*32)
	outE := os.Stdout

	svc, err := cc.GetService(ctx, name, opts...)
	if err != nil {
		return err
	}

	file, err := os.Open(pack)
	if err != nil {
		return err
	}

	s, err := cc.UpgradeService(ctx, name, v, opts...)
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

	fmt.Fprintf(outE, "upgrade service %s %s -> %s\n", name, svc.Version, v)
	return nil
}

func UpgradeServiceCmd() *cli.Command {
	return &cli.Command{
		Name:     "upgrade",
		Usage:    "upgrade a service",
		Category: "service",
		Action:   upgradeService,
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
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "specify the version for service",
			},
		},
	}
}
