package backup

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/backup/create"
	"github.com/skpr/cli/cmd/skpr/backup/list"
	"github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Manage the lifecycle for backups for an environment`
)

// NewCommand creates a new cobra.Command for 'backup' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "backup",
		DisableFlagsInUseLine: true,
		Short:                 "Backup application data. Databases, Files etc",
		Long:                  cmdLong,
		GroupID:               command.GroupDisasterRecovery,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
