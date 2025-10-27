package create

import (
	"github.com/spf13/cobra"

	v1create "github.com/skpr/cli/internal/command/volume/backup/create"
)

var (
	cmdLong = `Create a filesystem backup of an environment.`

	cmdExample = `
  # Create a filesystem backup of dev.
  skpr volume backup create dev

  # Create and wait for a filesystem backup.
  skpr volume backup create dev --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Create a filesystem backup",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for filesystem backup to complete")

	return cmd
}
