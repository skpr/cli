package set

import (
	"github.com/spf13/cobra"

	v1threshold "github.com/skpr/cli/internal/command/trace/set/threshold"
)

var (
	thresholdLong = `
  Set the tracing threshold for the specified environment.

  The threshold is the minimum duration a function call has to run for before it
  is traced. Compass only fires a probe for calls above it, so it is the main
  control over how much a traced environment pays for tracing. A threshold of
  zero traces every call, which is expensive on a busy environment.`

	thresholdExample = `
  # Set the tracing threshold for the dev environment
  skpr trace set threshold dev 10ms`
)

func newThresholdCommand() *cobra.Command {
	command := v1threshold.Command{}

	cmd := &cobra.Command{
		Use:                   "threshold <environment> <duration>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Set the tracing threshold for an environment",
		Long:                  thresholdLong,
		Example:               thresholdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Threshold = args[1]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
