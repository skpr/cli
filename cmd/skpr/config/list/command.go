package list

import (
	"github.com/spf13/cobra"

	v1list "github.com/skpr/cli/internal/command/config/list"
)

const (
	// MaxValueLength to be applied when listing values.
	MaxValueLength = 100
)

var (
	cmdLong = `List all of the config key/value pairs for the specified environment`

	cmdExample = `
  # List all of the config for dev environment
  skpr config list dev
	
  # List all of the config for dev environment as JSON
  skpr config list dev --json
	
  # List all of the config for dev environment with decoded secrets
  skpr config list dev --show-secrets`
)

// NewCommand creates a new cobra.Command for 'list' sub command
func NewCommand() *cobra.Command {
	command := v1list.Command{}

	cmd := &cobra.Command{
		Use:                   "list <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "List all of the config key/value pairs for the specified environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.JSON, "json", command.JSON, "Show output as JSON")
	cmd.Flags().BoolVar(&command.ShowSecrets, "show-secrets", command.ShowSecrets, "Show decoded secrets")
	cmd.Flags().StringVar(&command.FilterType, "filter-type", command.FilterType, "Filter based on config type")
	cmd.Flags().BoolVar(&command.Wide, "wide", command.Wide, "Show all config values")

	return cmd
}
