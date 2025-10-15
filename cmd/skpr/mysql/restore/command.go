package restore

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/mysql/restore/create"
	"github.com/skpr/cli/cmd/skpr/mysql/restore/list"
)

var (
	cmdLong = `Initiate a restore process which will restore a MySQL backup to a specified environment`

	cmdExample = `
  # Create a MySQL restore from a backup for dev environment.
  skpr mysql restore create dev BACKUP_ID
  
  # List MySQL restores which have been created for dev environment.
  skpr mysql restore list dev`
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
