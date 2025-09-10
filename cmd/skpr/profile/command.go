package profile

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/compass/cli/app"
	"github.com/spf13/cobra"
)

var (
	cmdLong = `
  Profile requests as they flow through your application using Compass.`

	cmdExample = `
  # Start profiling for the specified environment
  skpr profile <environment>`
)

// NewCommand creates a new cobra.Command for 'profile' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "profile [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Profile requests as they flow through your application",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(app.NewModel(""), tea.WithAltScreen())

			_, err := p.Run()
			if err != nil {
				return fmt.Errorf("failed to run program: %w", err)
			}

			return nil
		},
	}

	return cmd
}
