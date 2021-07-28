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
	"github.com/lack-io/cli"
	vclient "github.com/lack-io/vine/core/client"
)

func healthService(c *cli.Context) error {

	addr := c.String("host")
	cc := client.New(addr)

	ctx := context.Background()
	outE := os.Stdout

	err := cc.Healthz(ctx, vclient.WithAddress(addr))
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "OK\n")
	return nil
}

func HealthCmd() *cli.Command {
	return &cli.Command{
		Name:   "health",
		Usage:  "confirm gpmd status",
		Action: healthService,
	}
}