package create

import (
	"github.com/spf13/cobra"

	v1create "github.com/skpr/cli/internal/command/filesystem/backup/create"
)

var (
	cmdLong = `Create a filesystem backup of an environment.`

	cmdExample = `
  # Create a public filesystem backup of dev.
  skpr filesystem backup create dev public

  # Create and wait for a private filesystem backup.
  skpr filesystem backup create dev private --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create <environment> <id>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Create a filesystem backup",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.VolumeName = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for filesystem backup to complete")

	return cmd
}
