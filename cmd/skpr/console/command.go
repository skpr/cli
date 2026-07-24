package console

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	"github.com/skpr/cli/internal/command/console"
)

// NewCommand creates a new cobra.Command for 'console' sub command
func NewCommand() *cobra.Command {
	command := console.Command{}

	cmd := &cobra.Command{
		Use:                   "console <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Open the Skpr console for an environment in browser.",
		GroupID:               skprcommand.GroupDebug,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Print, "print", "p", false, "Only display the link instead of opening it")

	return cmd
}
