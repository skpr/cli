package pull

import (
	"github.com/spf13/cobra"

	v1pull "github.com/skpr/cli/internal/command/mysql/image/pull"
)

var (
	cmdLong = `
  Pulls a database image associated with an environment.`

	cmdExample = `
  # Pull the database image for an environment.
  skpr mysql image pull ENVIRONMENT`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1pull.Command{}

	cmd := &cobra.Command{
		Use:                   "pull",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Pull a database image for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Params.Environment = args[0]

			if len(args) > 1 {
				command.Params.Databases = args[1:]
			} else {
				command.Params.Databases = []string{"default"}
			}

			return command.Run(cmd.Context())
		},
	}

	return cmd
}
