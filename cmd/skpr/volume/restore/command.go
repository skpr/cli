package restore

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/volume/restore/create"
	"github.com/skpr/cli/cmd/skpr/volume/restore/list"
)

var (
	cmdLong = `Initiate a restore process which will restore a filesystem backup to a specified environment`

	cmdExample = `
  # Create a filesystem restore from a backup for dev environment.
  skpr volume restore create dev BACKUP_ID
  
  # List filesystem restores which have been created for dev environment.
  skpr volume restore list dev`
)

// NewCommand creates a new cobra.Command for 'restore' sub command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "restore",
		DisableFlagsInUseLine: true,
		Short:                 "Restore filesystem data.",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
