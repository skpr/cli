package create

import (
	v1create "github.com/skpr/cli/internal/command/mysql/restore/create"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Create a MySQL restore of an environment.`

	cmdExample = `
  # Create a Mysql restore for dev with BACKUP_ID and wait.
  skpr mysql restore create dev BACKUP_ID --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create [environment] [backup_id]",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Create a MySQL restore",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Backup = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Force, "force", "f", false, "Skpr will not request a confirmation for configuration changes.")
	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for restore to complete")

	return cmd
}
