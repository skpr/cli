package list

import (
	v1list "github.com/skpr/cli/internal/command/purge/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  List purge requests for a given Skpr environment.`

	cmdExample = `
  # List purge events that have been created on the dev environment.
  skpr mysql image list dev`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List purge requests",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
