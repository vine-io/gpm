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
	"time"

	json "github.com/json-iterator/go"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/client"
	"github.com/vine-io/pkg/unit"
	log "github.com/vine-io/vine/lib/logger"
	"gopkg.in/yaml.v3"
)

func infoService(c *cobra.Command, args []string) error {
	opts := getCallOptions(c)
	output, _ := c.Flags().GetString("output")
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	s, err := cc.Info(ctx, opts...)
	if err != nil {
		return err
	}

	switch output {
	case "wide":
		if runtime.GOOS == "windows" {
			log.Infof("'wide' don't support windows")
			b, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintf(outE, "%s\r\n", string(b))
			return nil
		}
		t := tw.NewWriter(os.Stdout)
		t.SetHeader([]string{"Property", "Value"})
		t.SetBorder(false)
		t.SetColWidth(200)
		t.SetAlignment(0)

		t.Append([]string{"Pid", fmt.Sprintf("# %d", s.Pid)})
		t.Append([]string{"Version", s.Version})
		t.Append([]string{"OS", s.Goos})
		t.Append([]string{"Arch", s.Arch})
		t.Append([]string{"Go version", s.Gov})
		if s.Stat != nil {
			t.Append([]string{"CPU", fmt.Sprintf("%.2f%%", s.Stat.CpuPercent)})
			t.Append([]string{"Memory", fmt.Sprintf("%s/%.1f%%", unit.ConvAuto(int64(s.Stat.Memory), 2), s.Stat.MemPercent)})
		}
		t.Append([]string{"UpTime", (time.Duration(s.UpTime) * time.Second).String()})
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

func InfoServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "get the information of gpmd",
		GroupID: "gpm",
		RunE:    infoService,
	}

	cmd.PersistentFlags().StringP("output", "o", "wide", "specify the format of output, example wide, json, yaml")

	return cmd
}
