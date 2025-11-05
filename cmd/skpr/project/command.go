package project

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/project/get"
	"github.com/skpr/cli/cmd/skpr/project/list"
	"github.com/skpr/cli/cmd/skpr/project/set"
	"github.com/skpr/cli/internal/command"
)

// NewCommand creates a new cobra.Command for 'config' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "project",
		DisableFlagsInUseLine: true,
		Short:                 "Manage projects",
		GroupID:               command.GroupLifecycle,
	}

	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(set.NewCommand())

	return cmd
}
