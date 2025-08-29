package create

import (
	"time"

	"github.com/spf13/cobra"

	v1create "github.com/skpr/cli/internal/command/v1/backup/create"
)

var (
	cmdLong = `
  Create a backup of an environment.`

	cmdExample = `
  # Create a backup.
  skpr backup create ENVIRONMENT

  # Create and wait for a backup.
  skpr backup create ENVIRONMENT --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Create a backup",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run()
		},
	}

	cmd.Flags().StringSliceVar(&command.MySQL.Policies, "mysql-policy", nil, "Policy to apply to MySQL Database Backup")
	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for backup to complete")
	cmd.Flags().DurationVar(&command.WaitTimeout, "wait-timeout", 4*time.Hour, "How long to wait for")
	cmd.Flags().Int32Var(&command.WaitErrorLimit, "wait-error-limit", 20, "How many errors to tolerate before exiting")

	return cmd
}
