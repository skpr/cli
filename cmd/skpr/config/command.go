package config

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/config/delete"
	"github.com/skpr/cli/cmd/skpr/config/get"
	"github.com/skpr/cli/cmd/skpr/config/list"
	"github.com/skpr/cli/cmd/skpr/config/set"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Manage application connection details, secrets, toggles and more.`
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		DisableFlagsInUseLine: true,
		Short:                 "Manage application connection details, secrets, toggles and more",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupLifecycle,
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(set.NewCommand())
	cmd.AddCommand(delete.NewCommand())

	return cmd
}
