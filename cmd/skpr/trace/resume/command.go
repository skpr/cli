package resume

import (
	"github.com/spf13/cobra"

	v1resume "github.com/skpr/cli/internal/command/trace/resume"
)

var (
	cmdLong = `Resume tracing for the specified environment.`

	cmdExample = `
  # Resume tracing for the dev environment
  skpr trace resume dev`
)

// NewCommand creates a new cobra.Command for 'resume' sub command
func NewCommand() *cobra.Command {
	command := v1resume.Command{}

	cmd := &cobra.Command{
		Use:                   "resume <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Resume tracing for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
