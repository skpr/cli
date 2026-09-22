package set

import (
	"github.com/spf13/cobra"

	v1set "github.com/skpr/cli/internal/command/trace/threshold/set"
)

var (
	cmdLong = `
  Set the tracing threshold for the specified environment.

  The threshold is the minimum duration a function call has to run for before it
  is traced. Compass only fires a probe for calls above it, so it is the main
  control over how much a traced environment pays for tracing. A threshold of
  zero traces every call, which is expensive on a busy environment.`

	cmdExample = `
  # Set the tracing threshold for the dev environment
  skpr trace threshold set dev 10ms`
)

// NewCommand creates a new cobra.Command for 'set' sub command
func NewCommand() *cobra.Command {
	command := v1set.Command{}

	cmd := &cobra.Command{
		Use:                   "set <environment> <duration>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Set the tracing threshold for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Threshold = args[1]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
