package trace

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1trace "github.com/skpr/cli/internal/command/trace"
)

var (
	cmdLong = `
  Trace requests as they flow through your application using Compass.`

	cmdExample = `
  # Start profiling for the specified environment
  skpr trace <environment>`
)

// NewCommand creates a new cobra.Command for 'trace' sub command
func NewCommand() *cobra.Command {
	command := v1trace.Command{}

	cmd := &cobra.Command{
		Use:                   "trace [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Trace requests as they flow through your application",
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
