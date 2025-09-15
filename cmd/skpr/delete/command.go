package delete

import (
	v1delete "github.com/skpr/cli/internal/command/delete"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Delete an environment.`

	cmdExample = `
  # Delete the dev environment
  skpr delete dev`
)

// NewCommand creates a new cobra.Command for 'delete' sub command
func NewCommand() *cobra.Command {
	command := v1delete.Command{}

	cmd := &cobra.Command{
		Use:                   "delete <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Delete a previously deployed environment",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Name = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(&command.DryRun, "dry-run", command.DryRun, "Displays list of resources that would be deleted.")
	cmd.Flags().BoolVar(&command.SkipConfirm, "yes", command.SkipConfirm, "Skips confirmation prompt (except if environment marked with \"production\")")
	cmd.Flags().BoolVarP(&command.Force, "force", "f", command.Force, "Skpr will not request a confirmation for configuration changes")
	_ = cmd.Flags().MarkHidden("force")

	return cmd
}
