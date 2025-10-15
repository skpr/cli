package job

import (
	"fmt"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(list.NewCommand())

	return cmd
}
