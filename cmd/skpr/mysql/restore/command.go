package restore

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/mysql/restore/create"
	"github.com/skpr/cli/cmd/skpr/mysql/restore/list"
)

var (
	cmdLong = `
  Initiate a restore process which will restore a MySQL backup to a specified environment`

	cmdExample = `
  # Create a MySQL restore from a backup.
  skpr mysql restore create ENVIRONMENT BACKUP_ID
  
  # List MySQL restores which have been created for an environment.
  skpr mysql restore list ENVIRONMENT`
)

// NewCommand creates a new cobra.Command for 'restore' sub command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "restore",
		DisableFlagsInUseLine: true,
		Short:                 "Restore MySQL data.",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
