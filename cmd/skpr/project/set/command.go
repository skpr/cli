package set

import (
	"github.com/spf13/cobra"

	v1set "github.com/skpr/cli/internal/command/project/set"
)

var (
	cmdExample = `
  # Set the contact for the project.
  skpr project set contact my-new-contact@example.com

  # Set the tags for the project.
  skpr project set tags "tag-a tag-b tag-c"`
)

// NewCommand creates a new cobra.Command for 'set' sub command
func NewCommand() *cobra.Command {
	command := v1set.Command{}

	cmd := &cobra.Command{
		Use:                   "set <key> <value>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Set an attribute for the current project.",
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Key = args[0]
			command.Value = args[1]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
