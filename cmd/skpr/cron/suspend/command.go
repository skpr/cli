package suspend

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/cron/suspend"
)

var (
	cmdLong = `Suspend all cron jobs for a given environment.`

	cmdExample = `
    # Suspend all cron jobs for dev environment
    skpr cron suspend dev`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := suspend.Command{}

	cmd := &cobra.Command{
		Use:                   "suspend [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Suspend all cron jobs associated with an environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Wait, "wait", "w", command.Wait, "Wait for running cron tasks to complete")
	cmd.Flags().DurationVarP(&command.Timeout, "timeout", "t", command.Timeout, "Allowed timeout threshold when waiting for cron tasks")

	return cmd
}
