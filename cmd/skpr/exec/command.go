package exec

import (
	v1exec "github.com/skpr/cli/internal/command/exec"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `Executes a shell command in the chosen environment.`

	cmdExample = `
  # Exec a shell command on the dev envirionment.
  skpr exec dev -- echo 'Hello!'`
)

// NewCommand creates a new cobra.Command for 'exec' sub command
func NewCommand() *cobra.Command {
	command := v1exec.Command{}

	cmd := &cobra.Command{
		Use:                   "exec <environment> -- <command>",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Execute a single shell command",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Command = args[1:]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
