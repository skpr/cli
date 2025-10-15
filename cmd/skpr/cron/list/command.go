package list

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/cron/list"
)

var (
	cmdLong = `List all the cron jobs for a given environment.`

	cmdExample = `
    # List all cron jobs for dev environment
    skpr cron list dev`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := list.Command{}

	cmd := &cobra.Command{
		Use:                   "list [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List all Crons associated with an environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
