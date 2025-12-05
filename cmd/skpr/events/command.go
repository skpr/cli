package events

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1events "github.com/skpr/cli/internal/command/events"
)

var (
	cmdLong = `List events for a specific environment.`

	cmdExample = `
  # List events for the dev environment
  skpr events dev`
)

// NewCommand creates a new cobra.Command for 'events' sub command
func NewCommand() *cobra.Command {
	command := v1events.Command{}

	cmd := &cobra.Command{
		Use:                   "events <environment>",
		Args:                  cobra.MatchAll(cobra.ExactArgs(1)),
		DisableFlagsInUseLine: true,
		Short:                 "List events for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupDebug,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
