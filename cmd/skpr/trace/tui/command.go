package tui

import (
	"github.com/spf13/cobra"

	v1tui "github.com/skpr/cli/internal/command/trace/tui"
)

var (
	cmdLong = `
  Watch requests as they flow through your application using Compass.

  Traces are streamed into an interactive application for as long as the command
  runs. Nothing is shown while tracing is suspended for the environment, or
  while no request runs for long enough to pass the tracing threshold.`

	cmdExample = `
  # Open the tracing TUI for the dev environment
  skpr trace tui dev`
)

// NewCommand creates a new cobra.Command for 'tui' sub command
func NewCommand() *cobra.Command {
	command := v1tui.Command{}

	cmd := &cobra.Command{
		Use:                   "tui <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Watch traces for an environment in an interactive application",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
