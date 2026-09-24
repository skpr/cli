package get

import (
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Show tracing settings for an environment.`

	cmdExample = `
  # Show the tracing threshold for the dev environment
  skpr trace get threshold dev`
)

// NewCommand creates a new cobra.Command for 'get' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		DisableFlagsInUseLine: true,
		Short:                 "Show tracing settings for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(newThresholdCommand())

	return cmd
}
