package suspend

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1suspend "github.com/skpr/cli/internal/command/suspend"
)

var (
	cmdLong = `Suspend an environment.`

	cmdExample = `
  # Suspend the dev environment
  skpr suspend dev`
)

// NewCommand creates a new cobra.Command for 'suspend' sub command
func NewCommand() *cobra.Command {
	command := v1suspend.Command{}

	cmd := &cobra.Command{
		Use:                   "suspend <environment>",
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Suspend an environment.",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupLifecycle,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
