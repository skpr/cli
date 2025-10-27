package get

import (
	v1get "github.com/skpr/cli/internal/command/project/get"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Get the contact for a project.`

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
		Short:                 "Get contact for a project.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
