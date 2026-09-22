package watch

import (
	"github.com/spf13/cobra"

	v1watch "github.com/skpr/cli/internal/command/trace/watch"
)

var (
	cmdLong = `
  Watch requests as they flow through your application using Compass.

  Traces are streamed into an interactive application for as long as the command
  runs. Nothing is shown while tracing is suspended for the environment, or
  while no request runs for long enough to pass the tracing threshold.`

	cmdExample = `
  # Watch traces for the dev environment
  skpr trace watch dev`
)

// NewCommand creates a new cobra.Command for 'watch' sub command
func NewCommand() *cobra.Command {
	command := v1watch.Command{}

	cmd := &cobra.Command{
		Use:                   "watch <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Watch traces for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
