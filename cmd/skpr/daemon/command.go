package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/daemon/list"
	"github.com/skpr/cli/cmd/skpr/daemon/resume"
	"github.com/skpr/cli/cmd/skpr/daemon/suspend"
)

var (
	cmdLong = `Daemon operations.`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "daemon",
		DisableFlagsInUseLine: true,
		Short:                 "Daemon operations.",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(resume.NewCommand())
	cmd.AddCommand(suspend.NewCommand())

	return cmd
}
