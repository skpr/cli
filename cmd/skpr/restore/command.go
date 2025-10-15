package restore

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/restore/create"
	"github.com/skpr/cli/cmd/skpr/restore/list"
)

var (
	cmdLong = "Initiate a restore process which will restore a backup to a specified environment"

	cmdExample = `
  # Create a restore from a backup.
  skpr restore create dev BACKUP_ID
  
  # List restores which have been created on the dev environment.
  skpr restore list dev`
)

// NewCommand creates a new cobra.Command for 'restore' sub command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "restore",
		DisableFlagsInUseLine: true,
		Short:                 "Restore application data. Databases, Files etc",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
