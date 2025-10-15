package set

import (
	v1set "github.com/skpr/cli/internal/command/config/set"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Set a config value for the specified environment.`

	cmdExample = `
  # Get auth.user config for dev environment
  skpr config set dev auth.user`
)

// NewCommand creates a new cobra.Command for 'set' sub command
func NewCommand() *cobra.Command {
	command := v1set.Command{}

	cmd := &cobra.Command{
		Use:                   "set [environment] [key] [value]",
		Args:                  cobra.ExactArgs(3),
		DisableFlagsInUseLine: true,
		Short:                 "Set a config value for the specified environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Key = args[1]
			command.Value = args[2]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.FromFile, "from-file", command.FromFile, "File which contains the value of key/value pair to set")
	cmd.Flags().BoolVar(&command.Secret, "secret", command.Secret, "Indicates the config value is a secret. This option encrypts\nthe value at rest, and obfuscates the value when displayed\nwith the config list command.")
	cmd.Flags().BoolVarP(&command.Force, "force", "f", command.Secret, "Skpr will not request a confirmation for configuration changes")

	return cmd
}
