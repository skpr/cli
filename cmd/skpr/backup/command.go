package backup

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/backup/create"
	"github.com/skpr/cli/cmd/skpr/backup/list"
)

var (
	cmdLong = `Manage the lifecycle for backups for an environment`
)

// NewCommand creates a new cobra.Command for 'backup' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "backup",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Backup application data. Databases, Files etc",
		Long:                  cmdLong,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
