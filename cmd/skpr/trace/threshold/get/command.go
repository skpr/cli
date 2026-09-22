package get

import (
	"github.com/spf13/cobra"

	v1get "github.com/skpr/cli/internal/command/trace/threshold/get"
)

var (
	cmdLong = `Show the tracing threshold for the specified environment.`

	cmdExample = `
  # Show the tracing threshold for the dev environment
  skpr trace threshold get dev`
)

// NewCommand creates a new cobra.Command for 'get' sub command
func NewCommand() *cobra.Command {
	command := v1get.Command{}

	cmd := &cobra.Command{
		Use:                   "get <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Show the tracing threshold for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
