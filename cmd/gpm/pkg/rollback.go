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

	"github.com/vine-io/gpm/pkg/runtime/client"
	"github.com/vine-io/cli"
)

func rollbackService(c *cli.Context) error {
	name := c.String("name")
	revision := c.String("revision")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if len(revision) == 0 {
		return fmt.Errorf("missing revision")
	}

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	s, err := cc.GetService(ctx, name, opts...)
	if err != nil {
		return err
	}

	err = cc.RollBackService(ctx, name, revision, opts...)
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "rollback %s %s -> %s\n", s.Name, s.Version, revision)
	return nil
}

func RollbackServiceCmd() *cli.Command {
	return &cli.Command{
		Name:     "rollback",
		Usage:    "reboot a service",
		Category: "service",
		Action:   rollbackService,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"N"},
				Usage:   "specify the name of service",
			},
			&cli.StringFlag{
				Name:    "revision",
				Aliases: []string{"R"},
				Usage:   "specify the revision of service",
			},
		},
	}
}
