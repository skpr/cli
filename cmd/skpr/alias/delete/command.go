package delete

import (
	"github.com/spf13/cobra"

	v1delete "github.com/skpr/cli/internal/command/alias/delete"
)

var (
	cmdLong = `
  Delete an alias of the command.`

	cmdExample = `
  # Delete an alias.
  skpr alias delete my-alias

  # Delete the alias and specify the skpr config directory.
  skpr alias delete my-alias --dir="/home/pnx/.skpr"`
)

// NewCommand creates a new cobra.Command for `skpr alias delete`.
func NewCommand() *cobra.Command {
	command := v1delete.Command{}

	cmd := &cobra.Command{
		Use:                   "delete",
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Delete your alias",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Alias = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.Dir, "dir", ".skpr", "The skpr config directory.")

	return cmd
}
