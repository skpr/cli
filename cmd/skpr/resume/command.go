package resume

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1resume "github.com/skpr/cli/internal/command/resume"
)

var (
	cmdLong = `Resume an environment.`

	cmdExample = `
  # Resume the dev environment
  skpr resume dev`
)

// NewCommand creates a new cobra.Command for 'resume' sub command
func NewCommand() *cobra.Command {
	command := v1resume.Command{}

	cmd := &cobra.Command{
		Use:                   "resume <environment>",
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Short:                 "Resume an environment.",
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
