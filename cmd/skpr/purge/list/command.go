package list

import (
	v1list "github.com/skpr/cli/internal/command/purge/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `List purge requests for a given Skpr environment.`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List purge requests",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
