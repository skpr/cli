package shell

import (
	v1shell "github.com/skpr/cli/internal/command/shell"
	"github.com/spf13/cobra"
)

var (
	cmdLong = "Connect to an environments command line session for running multiple commands."

	cmdExample = `
  # Connect to an environments command line session.
  skpr shell ENVIRONMENT`
)

// Options is the commandline options for 'shell' sub command
type Options struct{}

// NewOptions provides an instance of Options with default values
func NewOptions() Options {
	return Options{}
}

// NewCommand creates a new cobra.Command for 'shell' sub command
func NewCommand() *cobra.Command {
	command := v1shell.Command{}

	cmd := &cobra.Command{
		Use:                   "shell",
		DisableFlagsInUseLine: true,
		Short:                 "Execute a multiple shell commands in a session",
		Args:                  cobra.ExactArgs(1),
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
