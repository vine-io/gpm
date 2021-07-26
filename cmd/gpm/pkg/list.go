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
	"time"

	"github.com/gpm2/gpm/pkg/runtime"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/pkg/unit"
	"github.com/olekukonko/tablewriter"

	"github.com/lack-io/cli"
	"github.com/lack-io/vine/core/client"
	"github.com/lack-io/vine/core/client/grpc"
)

func listService(c *cli.Context) error {
	conn := grpc.NewClient(client.Retries(0))
	cc := pb.NewGpmService(runtime.GpmName, conn)

	addr := c.String("host")

	ctx := context.Background()

	in := &pb.ListServiceReq{}
	rsp, err := cc.ListService(ctx, in, client.WithAddress(addr))
	if err != nil {
		return err
	}

	tw := tablewriter.NewWriter(os.Stdout)
	tw.SetHeader([]string{"Name", "User", "Pid", "CPU", "Memory", "Status", "Uptime"})

	for _, item := range rsp.Services {
		row := make([]string, 0)
		row = append(row, item.Name)
		if item.SysProcAttr != nil {
			row = append(row, fmt.Sprintf("%s:%s", item.SysProcAttr.User, item.SysProcAttr.Group))
		} else {
			row = append(row, " ")
		}
		row = append(row, fmt.Sprintf("%d", item.Pid))
		row = append(row, fmt.Sprintf("%.1f%%", item.Stat.CpuPercent*100))
		row = append(row, fmt.Sprintf("%s/%.1f%%", unit.ConvAuto(int64(item.Stat.Memory), 1), item.Stat.MemPercent*100))
		row = append(row, item.Status)
		row = append(row, time.Now().Sub(time.Unix(item.StartTimestamp, 0)).String())
		tw.Append(row)
	}

	tw.SetColumnColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{},
		tablewriter.Colors{}, tablewriter.Colors{tablewriter.FgRedColor}, tablewriter.Colors{})
	tw.Render()
	fmt.Fprintf(os.Stdout, "\nTotal: %d\n", rsp.Total)

	return nil
}

func ListServicesCmd() *cli.Command {
	return &cli.Command{
		Name:     "list-service",
		Usage:    "list all local services",
		Category: "service",
		Action:   listService,
	}
}
