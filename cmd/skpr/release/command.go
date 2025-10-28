package release

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/release/info"
	"github.com/skpr/cli/cmd/skpr/release/list"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Find information on releases created from packaging your application.`

	cmdExample = `
  # List all releases.
  skpr release list

  # Show information on a release.
  skpr release info 1.0.0

  # Show information on a release in JSON format.
  skpr release info 1.0.0 --json`
)

// NewCommand creates a new cobra.Command for 'releases' sub command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "release",
		DisableFlagsInUseLine: true,
		Short:                 "Review releases which have been packaged for this project",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupLifecycle,
	}

	cmd.AddCommand(info.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
