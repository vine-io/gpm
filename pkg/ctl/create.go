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

func createService(c *cobra.Command, args []string) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout
	spec := &gpmv1.ServiceSpec{
		Env:         map[string]string{},
		SysProcAttr: &gpmv1.SysProcAttr{},
		Log:         &gpmv1.ProcLog{},
		AutoRestart: 1,
	}

	spec.Name, _ = c.Flags().GetString("name")
	spec.Bin, _ = c.Flags().GetString("bin")
	spec.Args, _ = c.Flags().GetStringSlice("args")
	spec.Dir, _ = c.Flags().GetString("dir")
	env, _ := c.Flags().GetStringSlice("env")
	spec.SysProcAttr.User, _ = c.Flags().GetString("user")
	spec.SysProcAttr.Group, _ = c.Flags().GetString("group")
	spec.Log.Expire, _ = c.Flags().GetInt32("log-expire")
	spec.Log.MaxSize, _ = c.Flags().GetInt64("log-max-size")
	spec.Version, _ = c.Flags().GetString("version")
	autoRestart, _ := c.Flags().GetBool("auto-restart")
	if err := spec.Validate(); err != nil {
		return err
	}
	for _, item := range env {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			spec.Env[parts[0]] = parts[1]
		}
	}
	if autoRestart {
		spec.AutoRestart = 1
	} else {
		spec.AutoRestart = 0
	}

	s, err := cc.CreateService(ctx, spec, opts...)
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "service '%s' %s\n", s.Name, s.Status)
	return nil
}

func CreateServiceCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create a service",
		GroupID: "service",
		RunE:    createService,
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
	cmd.PersistentFlags().Bool("auto-restart", true, "Whether auto restart service when it crashing")

	return cmd
}
