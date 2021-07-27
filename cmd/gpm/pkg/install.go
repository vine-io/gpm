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

	"github.com/gpm2/gpm/pkg/runtime/client"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/lack-io/cli"
	vclient "github.com/lack-io/vine/core/client"
)

func installService(c *cli.Context) error {

	addr := c.String("host")

	spec := &gpmv1.ServiceSpec{

	}

	spec.Name = c.String("name")
	if err := spec.Validate(); err != nil {
		return err
	}

	cc := client.New(addr)
	ctx := context.Background()
	outE := os.Stdout

	stream, err := cc.InstallService(ctx, spec, vclient.WithAddress(addr))
	if err != nil {
		return err
	}
	defer stream.Close()

	fmt.Fprintf(outE, "")
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
				Name:    "name",
				Aliases: []string{"N"},
				Usage:   "specify the name of service",
			},
		},
	}
}

//func install() {
//	app := vine.NewService()
//
//	cc := pb.NewGpmService(
//		runtime.GpmName, app.Client(),
//	)
//
//	ctx := context.Background()
//
//	in := &pb.InstallServiceReq{
//		Name:        "test",
//		Bin:         "/opt/web/web",
//		Args:        nil,
//		Dir:         "/opt/web",
//		Env:         nil,
//		Version:     "v1.0.0",
//		AutoRestart: 1,
//	}
//
//	rsp, err := cc.InstallService(ctx, client.WithRetries(0))
//	if err != nil {
//		log.Fatal(err)
//	}
//	f, err := os.Open("/tmp/web.tar.gz")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//	stat, _ := f.Stat()
//	pack := &gpmv1.Package{
//		Package: "/tmp/web.tar.gz",
//		Total:   stat.Size(),
//		Chunk:   nil,
//		IsOk:    false,
//	}
//
//	in.Pack = pack
//	go func() {
//		buf := make([]byte, 1024*32)
//		for {
//			n, err := f.Read(buf)
//			if err != nil && err != io.EOF {
//				log.Fatal(err)
//			}
//
//			if err == io.EOF {
//				in.Pack.IsOk = true
//			}
//
//			in.Pack.Chunk = buf[0:n]
//			rsp.Send(in)
//
//			if err == io.EOF {
//				break
//			}
//		}
//	}()
//
//	for {
//		out, err := rsp.Recv()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if out.Result.Error != "" {
//			log.Fatal(out.Result.Error)
//		}
//		if out.Result.IsOk {
//			break
//		}
//	}
//}
