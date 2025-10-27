package volume

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/volume/backup"
	"github.com/skpr/cli/cmd/skpr/volume/restore"
)

var (
	cmdLong = `Commands for interacting with the Skpr platforms filesystem features.`
)

// NewCommand creates a new cobra.Command for 'mysql' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "volume",
		DisableFlagsInUseLine: true,
		Short:                 "Perform filesystem tasks for an environment",
		Long:                  cmdLong,
	}

	cmd.AddCommand(backup.NewCommand())
	cmd.AddCommand(restore.NewCommand())

	return cmd
}
