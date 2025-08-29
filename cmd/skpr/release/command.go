package release

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/release/info"
	"github.com/skpr/cli/cmd/skpr/release/list"
)

var (
	cmdLong = `
  Find information on releases created from packaging your application.`

	cmdExample = `
  # List all releases.
  skpr release list

  # Show information on a release.
  skpr release info <release name>

  # Show information on a release in JSON format.
  skpr release info <release name> --json`
)

// NewCommand creates a new cobra.Command for 'releases' sub command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "release",
		DisableFlagsInUseLine: true,
		Short:                 "Review releases which have been packaged for this project",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(info.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
