package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/list"
)

var (
	cmdLong = `Overview of all environments and their current status`
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
