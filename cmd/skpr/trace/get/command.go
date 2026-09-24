package get

import (
	"fmt"

	"github.com/spf13/cobra"

	v1threshold "github.com/skpr/cli/internal/command/trace/get/threshold"
)

// SettingThreshold is the minimum duration a function call has to run for before it is traced.
const SettingThreshold = "threshold"

var (
	cmdLong = `
  Show a tracing setting for the specified environment.

  Settings:

    threshold  The minimum duration a function call has to run for before it
               is traced.`

	cmdExample = `
  # Show the tracing threshold for the dev environment
  skpr trace get dev threshold`
)

// NewCommand creates a new cobra.Command for 'get' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get <environment> <setting>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Show a tracing setting for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			environment, setting := args[0], args[1]

			switch setting {
			case SettingThreshold:
				command := v1threshold.Command{
					Environment: environment,
				}
				return command.Run(cmd.Context())
			default:
				return fmt.Errorf("unknown tracing setting %q, supported settings are: %s", setting, SettingThreshold)
			}
		},
	}

	return cmd
}
