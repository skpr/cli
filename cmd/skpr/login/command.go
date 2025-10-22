package login

import (
	"github.com/spf13/cobra"

	v1login "github.com/skpr/cli/internal/command/login"
)

var (
	cmdLong = "Login to the Skpr hosting platform."
)

// NewCommand creates a new cobra.Command for 'login' sub command
func NewCommand() *cobra.Command {
	command := v1login.Command{}

	cmd := &cobra.Command{
		Use:                   "login",
		DisableFlagsInUseLine: true,
		Short:                 "Login to the Skpr cluster.",
		Args:                  cobra.NoArgs,
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.Callback, "callback", "http://localhost:11218", "Endpoint to callback as a part of the OIDC workflow.")

	return cmd
}
