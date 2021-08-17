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
	"io"
	"os"
	"strings"

	"github.com/vine-io/cli"
	"google.golang.org/grpc/status"

	"github.com/vine-io/gpm/pkg/runtime/client"
	gpmv1 "github.com/vine-io/gpm/proto/apis/gpm/v1"
)

type Sender struct {
	in *gpmv1.TerminalIn
	s  *client.TerminalStream
}

func (s *Sender) Write(data []byte) (int, error) {
	n := len(data)
	s.in.Command = string(data)
	if err := s.s.Send(s.in); err != nil {
		return 0, err
	}
	return n, nil
}

func terminalBash(c *cli.Context) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()

	stdin := os.Stdin
	out := os.Stdout
	outE := os.Stderr

	in := &gpmv1.TerminalIn{}
	env := c.StringSlice("env")
	in.User = c.String("user")
	in.Group = c.String("group")
	if err := in.Validate(); err != nil {
		return err
	}
	for _, item := range env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			in.Env[parts[0]] = parts[1]
		}
	}

	t, err := cc.Terminal(ctx, opts...)
	if err != nil {
		return err
	}

	sender := &Sender{in: in, s: t}
	go io.Copy(sender, stdin)

	for {
		b, err := t.Recv()
		if err != nil {
			return errors.New(status.Convert(err).Message())
		}
		if b.Error != "" {
			fmt.Fprintf(outE, b.Error)
		}
		if b.Stderr != nil {
			fmt.Fprintf(outE, string(b.Stderr))
		}
		if b.Stdout != nil {
			fmt.Fprintf(out, string(b.Stdout))
		}
		if b.IsOk {
			break
		}
	}

	return nil
}

func TerminalBashCmd() *cli.Command {
	return &cli.Command{
		Name:     "terminal",
		Usage:    "start a terminal",
		Category: "bash",
		Action:   terminalBash,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"E"},
				Usage:   "specify the env for exec",
			},
			&cli.StringFlag{
				Name:  "user",
				Usage: "specify the user for exec",
			},
			&cli.StringFlag{
				Name:  "group",
				Usage: "specify the group for exec",
			},
		},
	}
}
