package shell

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1shell "github.com/skpr/cli/internal/command/shell"
)

var (
	cmdLong = "Connect to an environments command line session for running multiple commands."
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
		Use:                   "shell [environment]",
		DisableFlagsInUseLine: true,
		Short:                 "Execute a multiple shell commands in a session",
		Args:                  cobra.ExactArgs(1),
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupSecureShell,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	return cmd
}
