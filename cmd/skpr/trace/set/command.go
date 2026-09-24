package set

import (
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Change tracing settings for an environment.`

	cmdExample = `
  # Set the tracing threshold for the dev environment
  skpr trace set threshold dev 10ms`
)

// NewCommand creates a new cobra.Command for 'set' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set",
		DisableFlagsInUseLine: true,
		Short:                 "Change tracing settings for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(newThresholdCommand())

	return cmd
}
