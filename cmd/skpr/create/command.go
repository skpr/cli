package create

import (
	"github.com/spf13/cobra"

	v1create "github.com/skpr/cli/internal/command/v1/create"
)

var (
	cmdLong = `Create a new environment using a packaged release.`

	cmdExample = `
  # Create a dev environment from release 1.0.0
  skpr create dev 1.0.0`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create [environment] [release]",
		Args:                  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Create a new environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Version = args[1]
			return command.Run()
		},
	}

	return cmd
}
