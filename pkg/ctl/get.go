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
	"runtime"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/client"
	"github.com/vine-io/pkg/unit"
	"gopkg.in/yaml.v3"
)

func getService(c *cobra.Command, args []string) error {
	name, _ := c.Flags().GetString("name")
	output, _ := c.Flags().GetString("output")
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	s, err := cc.GetService(ctx, name, opts...)
	if err != nil {
		return err
	}

	switch output {
	case "wide":
		if runtime.GOOS == "windows" {
			b, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintf(outE, "%s\r\n", string(b))
		}
		t := tw.NewWriter(os.Stdout)
		t.SetHeader([]string{"Property", "Value"})
		t.SetBorder(false)
		t.SetColWidth(200)
		t.SetAlignment(0)

		t.Append([]string{"Name", s.Name})
		t.Append([]string{"Bin", s.Bin})
		t.Append([]string{"Args", strings.Join(s.Args, " ")})
		t.Append([]string{"Pid", fmt.Sprintf("# %d", s.Pid)})
		t.Append([]string{"Dir", s.Dir})
		env := make([]string, 0)
		for k, v := range s.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		t.Append([]string{"env", strings.Join(env, ",")})
		if s.SysProcAttr != nil {
			t.Append([]string{"User", fmt.Sprintf("user=%s, group=%s", s.SysProcAttr.User, s.SysProcAttr.Group)})
		}
		t.Append([]string{"Version", s.Version})
		if s.AutoRestart > 0 {
			t.Append([]string{"AutoRestart", "True"})
		} else {
			t.Append([]string{"AutoRestart", "False"})
		}
		if s.Stat != nil {
			t.Append([]string{"CPU", fmt.Sprintf("%.2f%%", s.Stat.CpuPercent)})
			t.Append([]string{"Memory", fmt.Sprintf("%s/%.1f%%", unit.ConvAuto(int64(s.Stat.Memory), 2), s.Stat.MemPercent)})
		}
		if s.Log != nil {
			t.Append([]string{"log expire", fmt.Sprintf("%d days", s.Log.Expire)})
			t.Append([]string{"log chunk", fmt.Sprintf("%s", unit.ConvAuto(s.Log.MaxSize, 2))})
		}
		t.Append([]string{"CreationTimestamp", time.Unix(s.CreationTimestamp, 0).String()})
		t.Append([]string{"UpdateTimestamp", time.Unix(s.UpdateTimestamp, 0).String()})
		t.Append([]string{"StartTimestamp", time.Unix(s.StartTimestamp, 0).String()})
		t.Append([]string{"Status", s.Status})
		if s.Msg != "" {
			t.Append([]string{"Message", s.Msg})
		}

		t.SetColumnColor(tw.Colors{tw.Bold}, tw.Colors{})
		t.Render()
	case "json":
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintf(outE, "%s\n", string(b))
	case "yaml":
		b, err := yaml.Marshal(s)
		if err != nil {
			return err
		}
		fmt.Fprintf(outE, "%s\n", string(b))
	default:
		return fmt.Errorf("invalid output value: %v", output)
	}

	return nil
}

func GetServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get service by name",
		GroupID: "service",
		RunE:    getService,
	}

	cmd.PersistentFlags().StringP("name", "N", "", "specify the name of service")
	cmd.PersistentFlags().StringP("output", "o", "wide", "specify the format of output, example wide, json, yaml")

	return cmd
}
