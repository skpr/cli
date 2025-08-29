package create

import (
	"github.com/spf13/cobra"

	v1create "github.com/skpr/cli/internal/command/v1/mysql/image/create"
)

var (
	cmdLong = `
  Create an image for a given Skpr environment.`

	cmdExample = `
  # Create an database image from the dev environment.
  skpr mysql image create ENVIRONMENT`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Create a database image from an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run()
		},
	}

	cmd.Flags().StringSliceVar(&command.Policies, "policy", command.Policies, "Name of the policy to apply to this image")
	cmd.Flags().StringVar(&command.Tag, "tag", command.Tag, "Tag to apply to this image. Will be prepended to database name eg. custom = dev-default-custom")

	return cmd
}
