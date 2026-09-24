package get

import (
	"github.com/spf13/cobra"

	v1threshold "github.com/skpr/cli/internal/command/trace/get/threshold"
)

var (
	thresholdLong = `
  Show the tracing threshold for the specified environment.

  The threshold is the minimum duration a function call has to run for before it
  is traced.`

	thresholdExample = `
  # Show the tracing threshold for the dev environment
  skpr trace get threshold dev`
)

func newThresholdCommand() *cobra.Command {
	command := v1threshold.Command{}

	cmd := &cobra.Command{
		Use:                   "threshold <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Show the tracing threshold for an environment",
		Long:                  thresholdLong,
		Example:               thresholdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
