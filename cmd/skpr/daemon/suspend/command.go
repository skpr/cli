package suspend

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/daemon/suspend"
)

var (
	cmdLong = `Suspend all daemon tasks for a given environment.`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := suspend.Command{}

	cmd := &cobra.Command{
		Use:                   "suspend <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Suspend all daemon tasks associated with an environment.",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Wait, "wait", "w", command.Wait, "Wait for running daemon tasks to complete.")
	cmd.Flags().DurationVarP(&command.Timeout, "timeout", "t", command.Timeout, "Allowed timeout threshold when waiting for daemon tasks.")

	return cmd
}
