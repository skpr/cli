package info

import (
	v1list "github.com/skpr/cli/internal/command/release/info"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Show information on an packaged release for this project.`

	cmdExample = `
  # Show information on release 1.0.0.
  skpr release info 1.0.0

  # Show information on release 1.0.0 in JSON format.
  skpr release info 1.0.0 --json`
)

// NewCommand creates a new cobra.Command for "info" subcommand.
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "info [version]",
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Show information on a release.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Version = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.JSON, "json", command.JSON, "Show output as JSON")

	return cmd
}
