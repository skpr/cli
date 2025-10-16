package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/alias/list"
)

var (
	cmdLong = `List all aliases.`

	cmdExample = `
  # List all aliases.
  skpr alias list`
)

// NewCommand creates a new cobra.Command for `skpr alias list`.
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.MatchAll(cobra.ExactArgs(0), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "List your aliases",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run()
		},
	}

	return cmd
}
