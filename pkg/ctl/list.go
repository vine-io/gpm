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
	"time"

	twr "github.com/olekukonko/tablewriter"
	"github.com/vine-io/cli"
	"github.com/vine-io/gpm/pkg/internal/client"
	"github.com/vine-io/pkg/unit"
)

func listService(c *cli.Context) error {

	opts := getCallOptions(c)
	cc := client.New()
	ctx := context.Background()
	outE := os.Stdout

	list, total, err := cc.ListService(ctx, opts...)
	if err != nil {
		return err
	}
	if total == 0 {
		return fmt.Errorf("no services")
	}

	tw := twr.NewWriter(outE)
	tw.SetHeader([]string{"Name", "Version", "User", "Pid", "CPU", "Memory", "Status", "Uptime"})

	for _, item := range list {
		row := make([]string, 0)
		row = append(row, item.Name)
		row = append(row, item.Version)
		if item.SysProcAttr != nil {
			row = append(row, fmt.Sprintf("%s:%s", item.SysProcAttr.User, item.SysProcAttr.Group))
		} else {
			row = append(row, " ")
		}
		row = append(row, fmt.Sprintf("%d", item.Pid))
		row = append(row, fmt.Sprintf("%.1f%%", item.Stat.CpuPercent))
		row = append(row, fmt.Sprintf("%s/%.1f%%", unit.ConvAuto(int64(item.Stat.Memory), 1), item.Stat.MemPercent))
		row = append(row, item.Status)
		row = append(row, time.Now().Sub(time.Unix(item.StartTimestamp, 0)).String())
		tw.Append(row)
	}

	tw.SetColumnColor(twr.Colors{}, twr.Colors{}, twr.Colors{}, twr.Colors{}, twr.Colors{},
		twr.Colors{}, twr.Colors{twr.FgRedColor}, twr.Colors{})
	tw.Render()
	fmt.Fprintf(os.Stdout, "\nTotal: %d\n", total)

	return nil
}

func ListServicesCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Usage:    "list all local services",
		Category: "service",
		Action:   listService,
	}
}
