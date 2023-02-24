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
	"time"

	"github.com/spf13/cobra"
	"github.com/vine-io/gpm/pkg/ctl"
	verrs "github.com/vine-io/vine/lib/errors"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gpm",
		Short:   "package manage tools",
		Version: ctl.GetGitTag(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid subcommand")
		},
	}

	rootCmd.ResetFlags()
	rootCmd.PersistentFlags().StringP("host", "H", "127.0.0.1:33700", "the ip address of gpmd")
	rootCmd.PersistentFlags().Duration("dial-timeout", time.Second*30, "specify dial timeout for call option")
	rootCmd.PersistentFlags().Duration("request-timeout", time.Second*30, "pecify request timeout for call option")

	rootCmd.ResetCommands()
	rootCmd.AddCommand(
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
	)

	//sort.Sort(cli.FlagsByName(app.Flags))
	//sort.Sort(cli.CommandsByName(app.Commands))

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stdout, "gpm exec: %s\n", verrs.FromErr(err).Detail)
	}
}
