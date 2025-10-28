package filesystem

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/filesystem/backup"
	"github.com/skpr/cli/cmd/skpr/filesystem/restore"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Commands for interacting with the Skpr platforms filesystem features.`
)

// NewCommand creates a new cobra.Command for 'mysql' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "filesystem",
		DisableFlagsInUseLine: true,
		Short:                 "Perform filesystem tasks for an environment",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupDataStorage,
	}

	cmd.AddCommand(backup.NewCommand())
	cmd.AddCommand(restore.NewCommand())

	return cmd
}
