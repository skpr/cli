package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/logs/list"
)

var (
	cmdLong = `List the streams of a running application.`

	cmdExample = `
  # List the streams of a running application.
  skpr logs list dev`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List the streams of a running application",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
