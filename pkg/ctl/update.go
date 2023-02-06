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
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	pbr "github.com/schollz/progressbar/v3"
	"github.com/vine-io/cli"
	"google.golang.org/grpc/status"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/client"
)

func update(c *cli.Context) error {

	pack := c.String("package")

	opts := getCallOptions(c)
	ctx := context.Background()
	outE := os.Stdout
	pb := pbr.NewOptions(0,
		pbr.OptionSetWriter(outE),
		pbr.OptionShowBytes(true),
		pbr.OptionEnableColorCodes(true),
		pbr.OptionOnCompletion(func() {
			fmt.Fprintf(outE, "\n")
		}),
	)

	cc := client.New()
	if pack == "" {
		pack = os.Args[0]
	}
	file, err := os.Open(pack)
	if err != nil {
		return err
	}
	defer file.Close()

	ech := make(chan error, 1)
	done := make(chan struct{}, 1)
	buf := make([]byte, 1024*32)

	stream, err := cc.Update(ctx, opts...)
	if err != nil {
		return err
	}

	p := &gpmv1.UpdateIn{Version: internal.GitTag}
	stat, _ := file.Stat()
	if stat != nil {
		p.Total = stat.Size()
	}
	pb.ChangeMax64(p.Total)

	go func() {
		for {
			b, err := stream.Recv()
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

	for {
		n, e := file.Read(buf)
		if e != nil && e != io.EOF {
			return e
		}

		if e == io.EOF {
			p.IsOk = true
		}

		_ = pb.Add(n)

		p.Length = int64(n)
		p.Chunk = buf[0:n]
		e1 := stream.Send(p)
		if e1 != nil {
			return e1
		}

		if e == io.EOF {
			break
		}
	}

	select {
	case e := <-ech:
		_ = pb.Clear()
		return e
	case <-done:
	}

	p.DeploySignal = true
	err = stream.Send(p)
	if err != nil {
		return fmt.Errorf("send start signal: %v", err)
	}

	time.Sleep(time.Second * 2)

	timeout := time.After(time.Minute)
	var n int
	for {
		select {
		case <-timeout:
			return fmt.Errorf("waitting for gpmd timeout")
		default:
		}
		cc = client.New()
		err := cc.Healthz(ctx, opts...)
		if err == nil && n > 3 {
			break
		}
		n += 1
		time.Sleep(time.Second * 2)
	}

	fmt.Fprintf(outE, "gpm update successfully!\n")
	return nil
}

func UpdateCmd() *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "update gpm and gpmd",
		Action: update,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"P"},
				Usage:   "specify the package for update",
			},
		},
	}
}
