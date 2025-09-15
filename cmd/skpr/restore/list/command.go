package list

import (
	v1list "github.com/skpr/cli/internal/command/restore/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
		List all of the available restores for a given Skpr environment.`

	cmdExample = `
		# List restores which have been created on the dev environment.
		skpr restore list ENVIRONMENT`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List restores for an environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.JSON, "json", "j", command.JSON, "Show output as JSON")

	return cmd
}
