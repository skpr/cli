package delete

import (
	v1delete "github.com/skpr/cli/internal/command/config/delete"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Delete a config value for the specified environment.`

	cmdExample = `
  # Get auth.user config for dev environment
  skpr config delete dev auth.user`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := v1delete.Command{}

	cmd := &cobra.Command{
		Use:                   "delete [environment] [key]",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Delete a configuration key/value pair",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Key = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Force, "force", "f", command.Force, "Skpr will not request a confirmation for configuration changes")

	return cmd
}
