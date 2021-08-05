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

package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gpm2/gpm/cmd/gpm/pkg"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/lack-io/cli"
	verrs "github.com/lack-io/vine/proto/apis/errors"
)

func main() {
	app := &cli.App{
		Name:    "gpm",
		Usage:   "package manage tools",
		Version: runtime.GetVersion(),
		Commands: []*cli.Command{
			pkg.HealthCmd(),
			pkg.DeployCmd(),
			pkg.TarCmd(),
			pkg.UpdateCmd(),
			pkg.RunCmd(),
			pkg.ShutdownCmd(),

			pkg.ListServicesCmd(),
			pkg.InfoServiceCmd(),
			pkg.GetServiceCmd(),
			pkg.CreateServiceCmd(),
			pkg.StartServiceCmd(),
			pkg.StopServiceCmd(),
			pkg.DeleteServiceCmd(),
			pkg.RebootServiceCmd(),
			pkg.TailServiceCmd(),

			pkg.InstallServiceCmd(),
			pkg.UpgradeServiceCmd(),
			pkg.RollbackServiceCmd(),
			pkg.VersionServiceCmd(),

			pkg.LsBashCmd(),
			pkg.ExecBashCmd(),
			pkg.PushBashCmd(),
			pkg.PullBashCmd(),
			pkg.TerminalBashCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Usage:   "the ip address of gpmd",
				EnvVars: []string{"GPM_HOST"},
				Value:   "127.0.0.1:7700",
			},
			&cli.DurationFlag{
				Name:    "dial-timeout",
				Usage:   "specify dial timeout for call option",
				EnvVars: []string{"GPM_DIAL_TIMEOUT"},
				Value:   time.Second * 30,
			},
			&cli.DurationFlag{
				Name:    "request-timeout",
				Usage:   "specify request timeout for call option",
				EnvVars: []string{"GPM_REQUEST_TIMEOUT"},
				Value:   time.Second * 30,
			},
		},
		EnableBashCompletion: true,
		Action: func(ctx *cli.Context) error {
			return fmt.Errorf("invalid subcommand")
		},
		Authors: []*cli.Author{{
			Name:  "lack",
			Email: "598223084@qq.com",
		}},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Version = runtime.GitTag
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stdout, "gpm exec: %s\n", verrs.FromErr(err).Detail)
	}
}
