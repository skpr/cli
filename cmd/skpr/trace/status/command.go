package status

import (
	"github.com/spf13/cobra"

	v1status "github.com/skpr/cli/internal/command/trace/status"
)

var (
	cmdLong = `
  Show the tracing status for the specified environment.

  This reports if tracing is suspended, along with the tracing threshold which
  is currently applied.`

	cmdExample = `
  # Show the tracing status for the dev environment
  skpr trace status dev`
)

// NewCommand creates a new cobra.Command for 'status' sub command
func NewCommand() *cobra.Command {
	command := v1status.Command{}

	cmd := &cobra.Command{
		Use:                   "status <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Show the tracing status for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
