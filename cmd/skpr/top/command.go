package top

import (
	"github.com/spf13/cobra"

	v1top "github.com/skpr/cli/internal/command/top"
)

var (
	cmdLong = `Display resource usage metrics for an environment.`

	cmdExample = `
  # Display CPU resource usage for the dev environment.
  skpr top dev`
)

// NewCommand creates a new cobra.Command for 'top' sub command
func NewCommand() *cobra.Command {
	command := v1top.Command{}

	cmd := &cobra.Command{
		Use:                   "top <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Display resource usage metrics for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
