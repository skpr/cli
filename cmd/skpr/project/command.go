package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/project/get"
	"github.com/skpr/cli/cmd/skpr/project/set"
)

var (
	cmdLong = `Manage a project.`

	cmdExample = `
	# Get details for the project.
	skpr project get`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "project",
		DisableFlagsInUseLine: true,
		Short:                 "Manage a project",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(set.NewCommand())

	return cmd
}
