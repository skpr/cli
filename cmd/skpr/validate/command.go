package validate

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/validate"
)

var (
	cmdLong = `Validate the configuration of a specific environment.`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "validate",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Validate the configuration of a specific environment.",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
