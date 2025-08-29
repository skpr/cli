package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/config/delete"
	"github.com/skpr/cli/cmd/skpr/config/get"
	"github.com/skpr/cli/cmd/skpr/config/list"
	"github.com/skpr/cli/cmd/skpr/config/set"
)

var (
	cmdLong = `Manage application connection details, secrets, toggles and more.`

	cmdExample = `
	# Get a single config value for the specified environment
	skpr config get [<flags>] <environment> <key>

	# List all of the config key/value pairs for the specified environment
	skpr config list [<flags>] <environment>

	# Set a config value for the specified environment
	skpr config set [<flags>] <environment> <key> [<value>]

	# Deletes a config value for the specified environment
	skpr config delete [<flags>] <environment> <key>`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
		// Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Manage application connection details, secrets, toggles and more",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no subcommand was provided, subcommand is required")
			}
			return nil
		},
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(set.NewCommand())
	cmd.AddCommand(delete.NewCommand())

	return cmd
}
