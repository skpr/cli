package threshold

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/trace/threshold/get"
	"github.com/skpr/cli/cmd/skpr/trace/threshold/set"
)

var (
	cmdLong = `
  Tracing threshold operations.

  The threshold is the minimum duration a function call has to run for before it
  is traced. Compass only fires a probe for calls above it, so it is the main
  control over how much a traced environment pays for tracing. A threshold of
  zero traces every call, which is expensive on a busy environment.`

	cmdExample = `
  # Show the tracing threshold for the dev environment
  skpr trace threshold get dev

  # Set the tracing threshold for the dev environment
  skpr trace threshold set dev 10ms`
)

// NewCommand creates a new cobra.Command for 'threshold' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "threshold",
		DisableFlagsInUseLine: true,
		Short:                 "Show, or set, the tracing threshold for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(set.NewCommand())

	return cmd
}
