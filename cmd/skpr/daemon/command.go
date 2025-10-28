package daemon

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/daemon/list"
	"github.com/skpr/cli/cmd/skpr/daemon/resume"
	"github.com/skpr/cli/cmd/skpr/daemon/suspend"
	skprcommand "github.com/skpr/cli/internal/command"
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
		GroupID:               skprcommand.GroupBackground,
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(resume.NewCommand())
	cmd.AddCommand(suspend.NewCommand())

	return cmd
}
