package alias

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/alias/delete"
	"github.com/skpr/cli/cmd/skpr/alias/list"
	"github.com/skpr/cli/cmd/skpr/alias/set"
)

var (
	cmdLong = `Manage aliases.`

	cmdExample = `
  # Create a new alias
  skpr alias set my-alias "echo 'Hello World'"

  # List all aliases
  skpr alias list

  # Delete an alias
  skpr alias delete my-alias

  # Create an alias in a different folder
  skpr alias set --dir="/path/to/.skpr" my-alias "echo 'Hello World'"

  # List all aliases in a different folder
  skpr alias list --dir="/path/to/.skpr"`
)

// NewCommand creates a new cobra.Command for 'alias' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "alias",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Create your own subcommands",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(set.NewCommand())

	return cmd
}
