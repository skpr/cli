package set

import (
	"fmt"

	"github.com/spf13/cobra"

	v1threshold "github.com/skpr/cli/internal/command/trace/set/threshold"
)

// SettingThreshold is the minimum duration a function call has to run for before it is traced.
const SettingThreshold = "threshold"

var (
	cmdLong = `
  Change a tracing setting for the specified environment.

  Settings:

    threshold  The minimum duration a function call has to run for before it
               is traced. Compass only fires a probe for calls above it, so it
               is the main control over how much a traced environment pays for
               tracing. A threshold of zero traces every call, which is
               expensive on a busy environment.`

	cmdExample = `
  # Set the tracing threshold for the dev environment
  skpr trace set dev threshold 10ms`
)

// NewCommand creates a new cobra.Command for 'set' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set <environment> <setting> <value>",
		Args:                  cobra.ExactArgs(3),
		DisableFlagsInUseLine: true,
		Short:                 "Change a tracing setting for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			environment, setting, value := args[0], args[1], args[2]

			switch setting {
			case SettingThreshold:
				command := v1threshold.Command{
					Environment: environment,
					Threshold:   value,
				}
				return command.Run(cmd.Context())
			default:
				return fmt.Errorf("unknown tracing setting %q, supported settings are: %s", setting, SettingThreshold)
			}
		},
	}

	return cmd
}
