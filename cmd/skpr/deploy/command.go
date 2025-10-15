package deploy

import (
	v1deploy "github.com/skpr/cli/internal/command/deploy"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Execute a deployment of a packaged release to a specified environment.`

	cmdExample = `
  # Deploy release 1.0.0 to the dev environment
  skpr deploy dev 1.0.0`
)

// NewCommand creates a new cobra.Command for 'deploy' sub command
func NewCommand() *cobra.Command {
	command := v1deploy.Command{}

	cmd := &cobra.Command{
		Use:                   "deploy [environment] [release]",
		Args:                  cobra.MatchAll(cobra.ExactArgs(2)),
		DisableFlagsInUseLine: true,
		Short:                 "Deploy a release to an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Version = args[1]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
