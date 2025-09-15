package logout

import (
	v1logout "github.com/skpr/cli/internal/command/logout"
	"github.com/spf13/cobra"
)

var (
	cmdLong = "Logout from the Skpr hosting platform."

	cmdExample = `
  # Logout from the Skpr hosting platform.
  skpr logout`
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
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
