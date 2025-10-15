package list

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/cron/job/list"
)

var (
	cmdLong = `List all the jobs that have been executed as part of cron.`

	cmdExample = `
    # List all jobs that have been executed for dev environment
    skpr cron job list dev`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := list.Command{}

	cmd := &cobra.Command{
		Use:                   "list [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List all jobs associated with an environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
