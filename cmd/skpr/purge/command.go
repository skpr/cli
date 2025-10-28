package purge

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/purge/create"
	"github.com/skpr/cli/cmd/skpr/purge/list"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `A series of commands for executing and reviewing cache invalidation.`

	cmdExample = `
  # Create a purge for a specific set of paths.
  skpr purge create dev /my-sub-path /my-sub-path-2
  
  # List purge requests for dev.
  skpr purge list dev`
)

// NewCommand creates a new cobra.Command for 'image' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "purge",
		DisableFlagsInUseLine: true,
		Short:                 "Perform MySQL tasks for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupCDN,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
