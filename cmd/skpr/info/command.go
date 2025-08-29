package info

import (
	"github.com/spf13/cobra"

	v1info "github.com/skpr/cli/internal/command/v1/info"
)

var (
	cmdLong = `
  Display information for Skpr environments.`

	cmdExample = `
  # Display information about the dev environment.
  skpr info dev`
)

// NewCommand creates a new cobra.Command for 'info' sub command
func NewCommand() *cobra.Command {
	command := v1info.Command{}

	cmd := &cobra.Command{
		Use:                   "info [environment]",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Get a detailed overview of an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Name = args[0]
			return command.Run()
		},
	}

	cmd.Flags().StringVarP(&command.Format, "format", "f", "json", "Output format - supported format is 'json'")

	return cmd
}
