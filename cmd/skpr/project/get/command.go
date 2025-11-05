package get

import (
	"github.com/spf13/cobra"

	v1get "github.com/skpr/cli/internal/command/project/get"
)

var (
	cmdExample = `
  # Get project details
  skpr project get`
)

// NewCommand creates a new cobra.Command for 'get' sub command
func NewCommand() *cobra.Command {
	command := v1get.Command{}

	cmd := &cobra.Command{
		Use:                   "get",
		DisableFlagsInUseLine: true,
		Short:                 "Get contact for the current project.",
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
