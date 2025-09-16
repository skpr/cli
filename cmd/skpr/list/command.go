package list

import (
	v1list "github.com/skpr/cli/internal/command/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
		Overview of all environments and their current status`

	cmdExample = `
		# List all environments
		skpr list`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Overview of all environments and their current status",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
