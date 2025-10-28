package dashboard

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1info "github.com/skpr/cli/internal/command/dashboard"
)

// NewCommand creates a new cobra.Command for 'info' sub command
func NewCommand() *cobra.Command {
	command := v1info.Command{}

	cmd := &cobra.Command{
		Use:                   "dashboard <environment>",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Open dashboards for an environment in browser.",
		GroupID:               skprcommand.GroupAuthentication,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Print, "print", "p", false, "Only display the link instead of opening it")

	return cmd
}
