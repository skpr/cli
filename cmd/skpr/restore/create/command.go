package create

import (
	v1create "github.com/skpr/cli/internal/command/restore/create"
	"time"

	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Create a restore of an environment.`

	cmdExample = `
  # Create a restore from a backup.
  skpr restore create ENVIRONMENT BACKUP_ID

  # Create and wait.
  skpr backup create ENVIRONMENT BACKUP_ID --wait`
)

// NewCommand creates a new cobra.Command for 'create' sub command
func NewCommand() *cobra.Command {
	command := v1create.Command{}

	cmd := &cobra.Command{
		Use:                   "create",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Create a restore",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Backup = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVarP(&command.Force, "force", "f", false, "Skpr will not request a confirmation for configuration changes.")
	cmd.Flags().BoolVar(&command.Wait, "wait", false, "Wait for restore to complete")
	cmd.Flags().DurationVar(&command.WaitTimeout, "wait-timeout", 4*time.Hour, "How long to wait for")
	cmd.Flags().Int32Var(&command.WaitErrorLimit, "wait-error-limit", 20, "How many errors to tolerate before exiting")

	return cmd
}
