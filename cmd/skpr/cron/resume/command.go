package resume

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/cron/resume"
)

var (
	cmdLong = `Resume all cron jobs for a given environment.`

	cmdExample = `
    # Resume all cron jobs for dev environment
    skpr cron resume dev`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := resume.Command{}

	cmd := &cobra.Command{
		Use:                   "resume [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Resume all cron jobs associated with an environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
