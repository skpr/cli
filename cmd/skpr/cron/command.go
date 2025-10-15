package cron

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/cron/job"
	"github.com/skpr/cli/cmd/skpr/cron/list"
	"github.com/skpr/cli/cmd/skpr/cron/resume"
	"github.com/skpr/cli/cmd/skpr/cron/suspend"
)

var (
	cmdLong = `Cron operations.`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "cron",
		DisableFlagsInUseLine: true,
		Short:                 "Cron operations.",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(job.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(resume.NewCommand())
	cmd.AddCommand(suspend.NewCommand())

	return cmd
}
