package pull

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/containers/docker"
	v1pull "github.com/skpr/cli/internal/command/release/pull"
)

var (
	cmdLong = `Pulls the packaged container images for a release.`

	cmdExample = `
  # Pull the packaged container images for a release.
  skpr release pull VERSION`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand(clientId docker.DockerClientId) *cobra.Command {
	command := v1pull.Command{}

	cmd := &cobra.Command{
		Use:                   "pull <version>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Pull release images.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Params.Name = args[0]
			command.ClientId = clientId

			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.Params.Service, "service", command.Params.Service, "A specific service image to pull")

	return cmd
}
