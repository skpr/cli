package mysql

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/mysql/image/create"
	"github.com/skpr/cli/cmd/skpr/mysql/image/list"
	"github.com/skpr/cli/cmd/skpr/mysql/image/pull"
)

var (
	cmdLong = `A series of commands for managing the lifecycle of mysql images on Skpr.`

	cmdExample = `
  # Create an database image from the dev environment.
  skpr mysql image create dev
  
  # List image which have been created on the dev environment.
  skpr mysql image list dev`
)

// NewCommand creates a new cobra.Command for 'image' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "image",
		DisableFlagsInUseLine: true,
		Short:                 "Perform MySQL tasks for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(pull.NewCommand())

	return cmd
}
