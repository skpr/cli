package cron

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/cron/job"
	"github.com/skpr/cli/cmd/skpr/cron/list"
	"github.com/skpr/cli/cmd/skpr/cron/resume"
	"github.com/skpr/cli/cmd/skpr/cron/suspend"
	skprcommand "github.com/skpr/cli/internal/command"
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
		GroupID:               skprcommand.GroupBackground,
	}

	cmd.AddCommand(job.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(resume.NewCommand())
	cmd.AddCommand(suspend.NewCommand())

	return cmd
}
