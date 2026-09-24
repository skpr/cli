package trace

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/trace/get"
	"github.com/skpr/cli/cmd/skpr/trace/resume"
	"github.com/skpr/cli/cmd/skpr/trace/set"
	"github.com/skpr/cli/cmd/skpr/trace/status"
	"github.com/skpr/cli/cmd/skpr/trace/suspend"
	"github.com/skpr/cli/cmd/skpr/trace/tui"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Trace requests as they flow through your application using Compass.`

	cmdExample = `
  # Watch traces for the dev environment in an interactive application
  skpr trace tui dev

  # Show the tracing status for the dev environment
  skpr trace status dev

  # Suspend and resume tracing for the dev environment
  skpr trace suspend dev
  skpr trace resume dev

  # Show, or set, the tracing threshold for the dev environment
  skpr trace get dev threshold
  skpr trace set dev threshold 10ms`
)

// NewCommand creates a new cobra.Command for 'trace' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "trace",
		DisableFlagsInUseLine: true,
		Short:                 "Trace requests as they flow through your application",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupDebug,
	}

	cmd.AddCommand(tui.NewCommand())
	cmd.AddCommand(status.NewCommand())
	cmd.AddCommand(suspend.NewCommand())
	cmd.AddCommand(resume.NewCommand())
	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(set.NewCommand())

	return cmd
}
