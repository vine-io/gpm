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

	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/internal/client"
	"google.golang.org/grpc/status"
)

func tailService(c *cobra.Command, args []string) error {
	opts := getCallOptions(c)
	name, _ := c.Flags().GetString("name")
	number, _ := c.Flags().GetInt64("number")
	follow, _ := c.Flags().GetBool("follow")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}

	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	s, err := cc.WatchServiceLog(ctx, name, number, follow, opts...)
	if err != nil {
		return err
	}

	for {
		b, err := s.Next()
		if err != nil && err != io.EOF {
			return errors.New(status.Convert(err).Message())
		}
		if err == io.EOF {
			s, err = cc.WatchServiceLog(ctx, name, number, follow, opts...)
			if err != nil {
				return err
			}
			continue
		}
		if b.Error != "" {
			return errors.New(b.Error)
		}
		fmt.Fprintln(outE, b.Text)
		if err == io.EOF {
			break
		}
	}

	return nil
}

func TailServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tail",
		Short:   "tail service logs",
		GroupID: "service",
		RunE:    tailService,
	}

	cmd.PersistentFlags().StringP("name", "N", "", "specify the name of service")
	cmd.PersistentFlags().Int64P("number", "n", 1024, "specify the number of service log")
	cmd.PersistentFlags().BoolP("follow", "f", false, "whether watching service log")

	return cmd
}
