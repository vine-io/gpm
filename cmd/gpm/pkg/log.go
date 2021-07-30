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
	"os"

	"github.com/gpm2/gpm/pkg/runtime/client"
	"github.com/lack-io/cli"
	"google.golang.org/grpc/status"
)

func tailService(c *cli.Context) error {
	opts := getCallOptions(c)
	name := c.String("name")
	number := c.Int64("number")
	follow := c.Bool("follow")
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
		if err != nil {
			return errors.New(status.Convert(err).Message())
		}
		if b.Error != "" {
			return errors.New(b.Error)
		}
		fmt.Fprintln(outE, b.Text)
	}

	return nil
}

func TailServiceCmd() *cli.Command {
	return &cli.Command{
		Name:     "tail",
		Usage:    "tail service logs",
		Category: "service",
		Action:   tailService,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"N"},
				Usage:   "specify the name of service",
			},
			&cli.Int64Flag{
				Name:    "number",
				Aliases: []string{"n"},
				Usage:   "specify the number of service log",
				Value:   1024,
			},
			&cli.BoolFlag{
				Name:    "follow",
				Aliases: []string{"f"},
				Usage:   "whether watching service log",
			},
		},
	}
}
