package logout

import (
	"github.com/spf13/cobra"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
