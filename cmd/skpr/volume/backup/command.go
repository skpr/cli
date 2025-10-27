package backup

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/volume/backup/create"
	"github.com/skpr/cli/cmd/skpr/volume/backup/list"
)

var (
	cmdLong = `Manage the lifecycle for filesystem backups for an environment`

	cmdExample = `
  # Create a backup for dev environment.
  skpr volume backup create dev

  # List all backups for dev environment.
  skpr volume backup list dev`
)

// NewCommand creates a new cobra.Command for 'backup' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "backup",
		DisableFlagsInUseLine: true,
		Short:                 "Backup filesystem data",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
