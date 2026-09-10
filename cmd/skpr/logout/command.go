package logout

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1logout "github.com/skpr/cli/internal/command/logout"
)

var (
	cmdLong = "Logout from the Skpr hosting platform."
)

// NewCommand creates a new cobra.Command for 'logout' sub command
func NewCommand() *cobra.Command {
	command := v1logout.Command{}

	cmd := &cobra.Command{
		Use:                   "logout",
		DisableFlagsInUseLine: true,
		Short:                 "Initiate a logout event from the Skpr hosting plstform",
		Args:                  cobra.NoArgs,
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupAuthentication,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	// Note: This is registered as a sign out URL with the identity provider, so
	// it needs to match exactly, including the path.
	cmd.Flags().StringVar(&command.Callback, "callback", "http://localhost:11218/logout", "Endpoint to callback as a part of the OIDC workflow.")

	return cmd
}
