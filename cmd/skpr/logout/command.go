package logout

import (
	"github.com/spf13/cobra"

	v1logout "github.com/skpr/cli/internal/command/v1/logout"
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
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Callback = args[0]
			return command.Run()
		},
	}

	return cmd
}
