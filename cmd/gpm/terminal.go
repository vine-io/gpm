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

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gpm2/gpm/pkg/runtime"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/vine"
	"github.com/lack-io/vine/core/client"
)

func main() {
	app := vine.NewService()

	cc := pb.NewGpmService(
		runtime.GpmName, app.Client(),
	)

	ctx := context.Background()

	in := &pb.TerminalReq{
		In: &gpmv1.TerminalIn{
			Command: "",
			Env:     map[string]string{"a": "bbb"},
			Uid:     0,
			Gid:     0,
		},
	}

	rsp, err := cc.Terminal(ctx, client.WithRetries(0))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		in.In.Command = `echo $a`
		e := rsp.Send(in)
		if e != nil {
			log.Fatal(err)
		}

		in.In.Command = `ifconfig`
		rsp.Send(in)

		in.In.Command = `exit`
		rsp.Send(in)
	}()

	for {
		out, err := rsp.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out.Result)
		if out.Result.IsOk {
			return
		}
	}
}
