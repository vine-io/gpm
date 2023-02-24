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

	"github.com/vine-io/cli"
	"github.com/vine-io/gpm/pkg/ctl"
	verrs "github.com/vine-io/vine/lib/errors"
)

func main() {
	app := &cli.App{
		Name:    "gpm",
		Usage:   "package manage tools",
		Version: ctl.GetVersion(),
		Commands: []*cli.Command{
			ctl.HealthCmd(),
			ctl.DeployCmd(),
			ctl.TarCmd(),
			ctl.UnTarCmd(),
			ctl.UpdateCmd(),
			ctl.RunCmd(),
			ctl.ShutdownCmd(),

			ctl.ListServicesCmd(),
			ctl.InfoServiceCmd(),
			ctl.GetServiceCmd(),
			ctl.CreateServiceCmd(),
			ctl.EditServiceCmd(),
			ctl.StartServiceCmd(),
			ctl.StopServiceCmd(),
			ctl.DeleteServiceCmd(),
			ctl.RestartServiceCmd(),
			ctl.TailServiceCmd(),

			ctl.InstallServiceCmd(),
			ctl.UpgradeServiceCmd(),
			ctl.RollbackServiceCmd(),
			ctl.ForgetServiceCmd(),
			ctl.VersionServiceCmd(),

			ctl.LsBashCmd(),
			ctl.ExecBashCmd(),
			ctl.PushBashCmd(),
			ctl.PullBashCmd(),
			ctl.TerminalBashCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Usage:   "the ip address of gpmd",
				EnvVars: []string{"GPM_HOST"},
				Value:   "127.0.0.1:33700",
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

	app.Version = ctl.GetGitTag()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stdout, "gpm exec: %s\n", verrs.FromErr(err).Detail)
	}
}
