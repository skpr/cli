package login

import (
	"github.com/spf13/cobra"

	v1login "github.com/skpr/cli/internal/command/v1/login"
)

var (
	cmdLong = "Login to the Skpr hosting platform."

	cmdExample = `
  # Login to the Skpr hosting platform.
  skpr login`
)

// NewCommand creates a new cobra.Command for 'login' sub command
func NewCommand() *cobra.Command {
	command := v1login.Command{}

	cmd := &cobra.Command{
		Use:                   "login",
		DisableFlagsInUseLine: true,
		Short:                 "Login to the Skpr cluster.",
		Args:                  cobra.ExactArgs(1),
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Callback = args[0]
			return command.Run()
		},
	}

	cmd.Flags().StringVar(&command.Callback, "callback", "http://localhost:11218", "Endpoint to callback as a part of the OIDC workflow.")

	return cmd
}
