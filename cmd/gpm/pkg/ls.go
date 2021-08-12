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
	"time"

	"github.com/vine-io/gpm/pkg/runtime/client"
	"github.com/lack-io/pkg/unit"
	"github.com/olekukonko/tablewriter"

	"github.com/vine-io/cli"
)

func lsBash(c *cli.Context) error {

	path := c.String("path")
	if path == "" {
		return fmt.Errorf("missing path")
	}

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	list, err := cc.Ls(ctx, path, opts...)
	if err != nil {
		return err
	}

	if len(list) > 0 {
		tw := tablewriter.NewWriter(outE)
		tw.SetHeader([]string{"Name", "Mode", "Size", "ModTime"})

		for _, item := range list {
			row := make([]string, 0)
			row = append(row, item.Name)
			row = append(row, item.Mode)
			row = append(row, unit.ConvAuto(item.Size, 2))
			row = append(row, time.Unix(item.ModTime, 0).String())
			tw.Append(row)
		}

		tw.Render()
	}

	return nil
}

func LsBashCmd() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Usage:    "list remote directory",
		Category: "bash",
		Action:   lsBash,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"P"},
				Usage:   "the specify the path for command",
			},
		},
	}
}
