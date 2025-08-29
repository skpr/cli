package backup

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/mysql/backup/create"
	"github.com/skpr/cli/cmd/skpr/mysql/backup/list"
)

var (
	cmdLong = `
  Manage the lifecycle for MySQL backups for an environment`

	cmdExample = `
  # Create a backup for an environment.
  skpr mysql backup create ENVIRONMENT

  # List all backups for an environment.
  skpr mysql backup list ENVIRONMENT`
)

// NewCommand creates a new cobra.Command for 'backup' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "backup",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Backup MySQL data",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
