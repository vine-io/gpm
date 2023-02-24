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
	"os"
	"strings"

	"github.com/spf13/cobra"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/client"
)

func editService(c *cobra.Command, args []string) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout
	spec := &gpmv1.EditServiceSpec{
		Env: map[string]string{},
	}

	name, _ := c.Flags().GetString("name")
	spec.Bin, _ = c.Flags().GetString("bin")
	spec.Args, _ = c.Flags().GetStringSlice("args")
	spec.Dir, _ = c.Flags().GetString("dir")
	env, _ := c.Flags().GetStringSlice("env")
	user, _ := c.Flags().GetString("user")
	group, _ := c.Flags().GetString("group")
	if user != "" || group != "" {
		spec.SysProcAttr = &gpmv1.SysProcAttr{}
	}
	if user != "" {
		spec.SysProcAttr.User = user
	}
	if group != "" {
		spec.SysProcAttr.Group = group
	}

	expire, _ := c.Flags().GetInt32("log-expire")
	maxSize, _ := c.Flags().GetInt64("log-max-size")
	if expire > 0 || maxSize > 0 {
		spec.Log = &gpmv1.ProcLog{}
	}
	if expire > 0 {
		spec.Log.Expire = expire
	}
	if maxSize > 0 {
		spec.Log.MaxSize = maxSize
	}

	spec.AutoRestart, _ = c.Flags().GetInt32("auto-restart")
	if err := spec.Validate(); err != nil {
		return err
	}
	for _, item := range env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			spec.Env[parts[0]] = parts[1]
		}
	}

	s, err := cc.EditService(ctx, name, spec, opts...)
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "edit service '%s' successfully!\n", s.Name)
	return nil
}

func EditServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "update a service parameters",
		GroupID: "service",
		RunE:    editService,
	}

	cmd.PersistentFlags().StringP("name", "N", "", "specify the name for service")
	cmd.PersistentFlags().StringP("bin", "B", "", "specify the bin for service")
	cmd.PersistentFlags().StringSliceP("args", "A", []string{}, "specify the args for service")
	cmd.PersistentFlags().StringP("dir", "D", "", "specify the root directory for service")
	cmd.PersistentFlags().StringP("env", "E", "", "specify the env for service")
	cmd.PersistentFlags().String("user", "", "specify the user for service")
	cmd.PersistentFlags().String("group", "", "specify the group for service")
	cmd.PersistentFlags().Int("log-expire", 15, "specify the expire for service log")
	cmd.PersistentFlags().Int64("log-max-size", 1024*1024*10, "specify the max size for service log")
	cmd.PersistentFlags().StringP("version", "V", "", "specify the version for service")
	cmd.PersistentFlags().Int("auto-restart", 1, "Whether auto restart service when it crashing")

	return cmd
}
