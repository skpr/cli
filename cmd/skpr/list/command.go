package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/v1/list"
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := command.Run(); err != nil {
				panic(err)
			}
		},
	}

	return cmd
}
