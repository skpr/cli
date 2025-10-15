package list

import (
	v1list "github.com/skpr/cli/internal/command/mysql/restore/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `List all of the available MySQL restores for a given Skpr environment.`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List MySQL restores for an environment",
		Long:                  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.JSON, "json", "j", command.JSON, "Show output as JSON")

	return cmd
}
