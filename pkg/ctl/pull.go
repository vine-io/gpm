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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal/client"

	pbr "github.com/schollz/progressbar/v3"
	"github.com/vine-io/pkg/unit"
	"google.golang.org/grpc/status"
)

func pullBash(c *cobra.Command, args []string) error {

	src, _ := c.Flags().GetString("src")
	regular, _ := c.Flags().GetBool("regular")
	dst, _ := c.Flags().GetString("dst")
	if src == "" {
		return fmt.Errorf("missing src")
	}
	if dst == "" {
		return fmt.Errorf("missing dst")
	}

	dir := dst
	if !regular {
		dir = filepath.Dir(dst)
	}
	stat, _ := os.Stat(dir)
	if stat == nil {
		return fmt.Errorf("invalid dst, local directory %s not exists", dir)
	}

	opts := getCallOptions(c)
	ctx := context.Background()
	outE := os.Stdout
	cc := client.New()
	pb := pbr.NewOptions(0,
		pbr.OptionSetWriter(outE),
		pbr.OptionShowBytes(true),
		pbr.OptionEnableColorCodes(true),
		pbr.OptionOnCompletion(func() {
			fmt.Fprintf(outE, "\n")
		}),
	)

	stream, err := cc.Pull(ctx, src, regular, opts...)
	if err != nil {
		return err
	}

	var name string
	var file *os.File
	for {
		var b *gpmv1.PullResult
		b, err = stream.Next()
		if err != nil {
			err = errors.New(status.Convert(err).Message())
			break
		}
		if b.Error != "" {
			err = errors.New(b.Error)
			break
		}
		if b.Name != name && b.Name != "" {
			name = b.Name
			pb.Reset()
			if file != nil {
				_ = file.Close()
			}
			pb.ChangeMax64(b.Total)
			pbn := b.Name
			if len(pbn) > 16 {
				pbn = "..." + pbn[len(pbn)-13:]
			}
			pb.Describe(fmt.Sprintf("download [%16s] [total: %8s]", pbn, unit.ConvAuto(b.Total, 2)))

			if regular {
				t := strings.TrimPrefix(b.Name, src)
				if strings.Contains(t, "/") || strings.Contains(t, "\\") {
					t = strings.ReplaceAll(t, "\\", "/")
					newDir := filepath.Join(dst, filepath.Dir(t))
					_ = os.MkdirAll(newDir, 0o777)
				}
				file, err = os.Create(filepath.Join(dst, t))
			} else {
				file, err = os.Create(dst)
			}
			if err != nil {
				break
			}
		}

		if b.Length > 0 {
			var n int
			n, err = file.Write(b.Chunk[0:b.Length])
			if err != nil {
				err = fmt.Errorf("write %s failed: %v", file.Name(), err)
				break
			}

			_ = pb.Add(n)
		}

		if b.Finished {
			err = nil
			break
		}
	}
	if file != nil {
		_ = file.Close()
	}
	if err != nil {
		_ = pb.Clear()
		return err
	}

	return nil
}

func PullBashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull",
		Short:   "pull file from service",
		GroupID: "bash",
		RunE:    pullBash,
	}

	cmd.PersistentFlags().StringP("src", "S", "", "specify the source for pull")
	cmd.PersistentFlags().StringP("dst", "D", "", "specify the target for pull")
	cmd.PersistentFlags().BoolP("regular", "R", false, "whether pull all files")

	return cmd
}
