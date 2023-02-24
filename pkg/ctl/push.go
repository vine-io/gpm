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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	pbr "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	client2 "github.com/vine-io/gpm/pkg/client"
	vclient "github.com/vine-io/vine/core/client"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
)

func pushBash(c *cobra.Command, args []string) error {

	src, _ := c.Flags().GetString("src")
	dst, _ := c.Flags().GetString("dst")
	if src == "" {
		return fmt.Errorf("missing src")
	}
	if dst == "" {
		return fmt.Errorf("missing dst")
	}

	stat, _ := os.Stat(src)
	if stat == nil {
		return fmt.Errorf("src %s not exists", src)
	}

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

	var err error
	if stat.IsDir() {
		err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if src == path || d.IsDir() {
				return nil
			}

			ddst := filepath.Join(dst, strings.TrimPrefix(path, src))
			pb.Describe(fmt.Sprintf("push [%30s]", path))
			pb.Reset()
			return push(ctx, pb, path, ddst, opts...)
		})
	} else {
		err = push(ctx, pb, src, dst, opts...)
	}

	return err
}

func push(ctx context.Context, pb *pbr.ProgressBar, src, dst string, opts ...vclient.CallOption) error {
	cc := client2.New()
	buf := make([]byte, 1024*32)

	stream, err := cc.Push(ctx, opts...)
	if err != nil {
		return err
	}

	err = send(src, dst, pb, stream, buf)
	if err != nil {
		return fmt.Errorf("send data: %v", err)
	}

	return stream.Wait()
}

func send(path, dst string, bar *pbr.ProgressBar, stream *client2.PushStream, buf []byte) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	p := &gpmv1.PushIn{Dst: dst, Name: filepath.Base(path)}
	stat, _ := file.Stat()
	if stat != nil {
		p.Total = stat.Size()
	}
	bar.ChangeMax64(p.Total)
	for {
		n, e := file.Read(buf)
		if e != nil && e != io.EOF {
			return e
		}

		_ = bar.Add(n)

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

	return nil
}

func PushBashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "push",
		Short:   "push files",
		GroupID: "bash",
		RunE:    pushBash,
	}

	cmd.PersistentFlags().StringP("src", "S", "", "specify the source for pull")
	cmd.PersistentFlags().StringP("dst", "D", "", "specify the target for pull")

	return cmd
}
