package set

import (
	"github.com/spf13/cobra"

	v1set "github.com/skpr/cli/internal/command/v1/alias/set"
)

var (
	cmdLong = `
  Set an alias for a command.`

	cmdExample = `
  # Set an alias.
  skpr alias set [<flags>] <alias> <expansion>

  # Set an alias
  skpr alias set my-alias "echo 'Hello World'"

  # Set the alias and specify the skpr config directory.
  skpr alias set --dir="/path/to/.skpr" <alias> <expansion>`
)

// NewCommand creates a new cobra.Command for `skpr alias set`.
func NewCommand() *cobra.Command {
	command := v1set.Command{}

	cmd := &cobra.Command{
		Use:                   "set",
		Args:                  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Set your alias",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Alias = args[0]
			command.Expansion = args[1]
			return command.Run()
		},
	}

	cmd.Flags().StringVar(&command.Dir, "dir", ".skpr", "The skpr config directory.")

	return cmd
}
