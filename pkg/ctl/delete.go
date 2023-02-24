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

	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/client"
)

func deleteService(c *cobra.Command, args []string) error {

	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout
	opts := getCallOptions(c)

	name, _ := c.Flags().GetString("name")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}

	s, err := cc.DeleteService(ctx, name, opts...)
	if err != nil {
		return err
	}

	fmt.Fprintf(outE, "service '%s' deleted\n", s.Name)
	return nil
}

func DeleteServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete a service",
		GroupID: "service", RunE: deleteService,
	}

	cmd.PersistentFlags().StringP("name", "N", "", "specify the name of service")

	return cmd
}
