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
	"io"
	"log"
	"os"

	"github.com/gpm2/gpm/pkg/runtime"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/vine"
	"github.com/lack-io/vine/core/client"
)

func push() {
	app := vine.NewService()

	cc := pb.NewGpmService(
		runtime.GpmName, app.Client(),
	)

	ctx := context.Background()

	in := &gpmv1.PushIn{
		Dst:   "/tmp/web.bbb",
		Chunk: nil,
		IsOk:  false,
	}

	rsp, err := cc.Push(ctx, client.WithRetries(0))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("/tmp/web.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			fmt.Println(n)
			if n > 0 {
				in.Chunk = buf[0:n]
				rsp.Send(&pb.PushReq{In: in})
			}
			if err == io.EOF {
				in.IsOk = true
				rsp.Send(&pb.PushReq{In: in})
				return
			}
		}
	}()

	for {
		result, err := rsp.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		if result.Result.IsOk {
			return
		}
	}
}
