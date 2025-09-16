package create

import (
	v1create "github.com/skpr/cli/internal/command/mysql/backup/create"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Create a MySQL backup of an environment.`

	cmdExample = `
  # Create a MySQL backup.
  skpr mysql backup create ENVIRONMENT

  # Create and wait for a MySQL backup.
  skpr mysql backup create ENVIRONMENT --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Create a MySQL backup",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for MySQL backup to complete")

	return cmd
}
