package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/release/list"
)

var (
	cmdLong = `List releases which have been packaged for this project.`

	cmdExample = `
  # List all releases.
  skpr release list

  # List all releases in JSON format.
  skpr release list --json`
)

// NewCommand creates a new cobra.Command for "list" subcommand.
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.MatchAll(cobra.ExactArgs(0)),
		DisableFlagsInUseLine: true,
		Short:                 "List all releases",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.Params.JSON, "json", command.Params.JSON, "Show output as JSON")

	return cmd
}
