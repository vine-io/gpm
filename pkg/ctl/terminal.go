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
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal/client"
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

func terminalBash(c *cobra.Command, args []string) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()

	stdin := os.Stdin
	out := os.Stdout
	outE := os.Stderr

	in := &gpmv1.TerminalIn{}
	env, _ := c.Flags().GetStringSlice("env")
	in.User, _ = c.Flags().GetString("user")
	in.Group, _ = c.Flags().GetString("group")
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
		if text := string(b.Error); text != "" {
			fmt.Fprintf(outE, text)
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

func TerminalBashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "terminal",
		Short:   "start a terminal",
		GroupID: "bash",
		RunE:    terminalBash,
	}

	cmd.PersistentFlags().StringSliceP("env", "E", []string{}, "specify the env for exec")
	cmd.PersistentFlags().String("user", "", "specify the user for exec")
	cmd.PersistentFlags().String("group", "", "specify the group for exec")

	return cmd
}
