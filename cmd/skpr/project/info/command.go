package info

import (
	"github.com/spf13/cobra"

	v1get "github.com/skpr/cli/internal/command/project/info"
)

var (
	cmdExample = `
  # Get project information
  skpr project info`
)

// NewCommand creates a new cobra.Command for 'info' sub command
func NewCommand() *cobra.Command {
	command := v1get.Command{}

	cmd := &cobra.Command{
		Use:                   "info",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Get full details for the current project.",
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
