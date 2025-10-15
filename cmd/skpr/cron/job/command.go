package job

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/cron/job/list"
)

var (
	cmdLong = `Job operations.`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "job",
		DisableFlagsInUseLine: true,
		Short:                 "Job operations.",
		Long:                  cmdLong,
	}

	cmd.AddCommand(list.NewCommand())

	return cmd
}
