package get

import (
	v1get "github.com/skpr/cli/internal/command/config/get"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Get a config value for the specified environment.`

	cmdExample = `
  # Get auth.user config for dev environment
  skpr config get dev auth.user`
)

// NewCommand creates a new cobra.Command for 'get' sub command
func NewCommand() *cobra.Command {
	command := v1get.Command{}

	cmd := &cobra.Command{
		Use:                   "get [environment] [key]",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Get a configuration key/value pair",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Key = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.ShowSecrets, "show-secret", true, "Show decoded secrets")

	return cmd
}
