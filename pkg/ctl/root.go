package ctl

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func ExecCmd() error {
	rootCmd := &cobra.Command{
		Use:           "gpm",
		Short:         "package manage tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       GetGitTag(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid subcommand")
		},
	}

	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "service",
			Title: "Service Subcommands",
		},
		&cobra.Group{
			ID:    "bash",
			Title: "Bash Subcommands",
		},
		&cobra.Group{
			ID:    "gpm",
			Title: "Gpm Subcommands",
		},
	)

	rootCmd.ResetCommands()
	rootCmd.AddCommand(
		HealthCmd(),
		DeployCmd(),
		TarCmd(),
		UnTarCmd(),
		UpdateCmd(),
		ShutdownCmd(),

		ListServicesCmd(),
		InfoServiceCmd(),
		GetServiceCmd(),
		CreateServiceCmd(),
		EditServiceCmd(),
		StartServiceCmd(),
		StopServiceCmd(),
		DeleteServiceCmd(),
		RestartServiceCmd(),
		TailServiceCmd(),

		InstallServiceCmd(),
		UpgradeServiceCmd(),
		RollbackServiceCmd(),
		ForgetServiceCmd(),
		VersionServiceCmd(),

		LsBashCmd(),
		ExecBashCmd(),
		PushBashCmd(),
		PullBashCmd(),
		TerminalBashCmd(),
	)

	runCmd, err := RunCmd()
	if err != nil {
		return err
	}
	rootCmd.AddCommand(runCmd)

	rootCmd.ResetFlags()
	rootCmd.PersistentFlags().StringP("host", "H", "127.0.0.1:33700", "the ip address of gpmd")
	rootCmd.PersistentFlags().Duration("dial-timeout", time.Second*30, "specify dial timeout for call option")
	rootCmd.PersistentFlags().Duration("request-timeout", time.Second*30, "pecify request timeout for call option")

	return rootCmd.Execute()
}
