package create

import (
	v1create "github.com/skpr/cli/internal/command/purge/create"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Create a purge request to invalidate edge caching.`

	cmdExample = `
  # Create a purge request.
  skpr purge create ENVIRONMENT PATH PATH PATH`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Create a purge request",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Paths = args[1:]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
