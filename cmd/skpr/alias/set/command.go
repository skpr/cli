package set

import (
	"strings"

	"github.com/spf13/cobra"

	v1set "github.com/skpr/cli/internal/command/alias/set"
)

var (
	cmdLong = `Set an alias for a command.`

	cmdExample = `
  # Set an alias
  skpr alias set my-alias "echo 'Hello World'"

  # Set the alias and specify the skpr config directory.
  alias set my-alias "echo 'Hello World'" --dir="/path/to/.skpr"`
)

// NewCommand creates a new cobra.Command for `skpr alias set`.
func NewCommand() *cobra.Command {
	command := v1set.Command{}

	cmd := &cobra.Command{
		Use:                   "set <alias> <expansion>",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Set your alias",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Alias = args[0]
			command.Expansion = strings.Join(args[1:], " ")
			return command.Run()
		},
	}

	return cmd
}
