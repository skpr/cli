package mysql

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/mysql/backup"
	img "github.com/skpr/cli/cmd/skpr/mysql/image"
	"github.com/skpr/cli/cmd/skpr/mysql/restore"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Commands for interacting with the Skpr platforms MySQL features.`
)

// NewCommand creates a new cobra.Command for 'mysql' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "mysql",
		DisableFlagsInUseLine: true,
		Short:                 "Perform MySQL tasks for an environment",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupDataStorage,
	}

	cmd.AddCommand(img.NewCommand())
	cmd.AddCommand(backup.NewCommand())
	cmd.AddCommand(restore.NewCommand())

	return cmd
}
