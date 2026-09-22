package suspend

import (
	"github.com/spf13/cobra"

	v1suspend "github.com/skpr/cli/internal/command/trace/suspend"
)

var (
	cmdLong = `Suspend tracing for the specified environment.`

	cmdExample = `
  # Suspend tracing for the dev environment
  skpr trace suspend dev`
)

// NewCommand creates a new cobra.Command for 'suspend' sub command
func NewCommand() *cobra.Command {
	command := v1suspend.Command{}

	cmd := &cobra.Command{
		Use:                   "suspend <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Suspend tracing for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
