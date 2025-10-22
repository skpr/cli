package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/mysql/backup/list"
)

var (
	cmdLong = `List all MySQL backups for an environment.`

	cmdExample = `
  # List MySQL backups for dev environment.
  skpr mysql backup list dev

  # List all MySQL backups for dev environment in JSON format.
  skpr mysql backup list dev --json

  # Pipe a list of all MySQL backups to jq for advanced query functionality.
  skpr mysql backup list dev --json | jq`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List MySQL backups for an environment",
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
