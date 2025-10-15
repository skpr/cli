package list

import (
	v1list "github.com/skpr/cli/internal/command/backup/list"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `List all backups for an environment.`

	cmdExample = `
  # List all backup for dev environment in JSON format.
  skpr backup list dev --json

  # Pipe a list of all backups to jq for advanced query functionality.
  skpr backup list dev --json | jq`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List backups for an environment",
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
