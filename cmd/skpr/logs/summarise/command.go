package summarise

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/logs/filter"
	v1summarise "github.com/skpr/cli/internal/command/logs/summarise"
)

var (
	cmdLong = `Summarise a window of an application's logs using a natural-language prompt.`

	cmdExample = `
  # Summarise the last hour of the default streams.
  skpr logs summarise dev

  # Summarise specific streams over the last 24 hours.
  skpr logs summarise dev nginx fpm --timeframe 24h

  # Ask a specific question of the summariser.
  skpr logs summarise dev --prompt "What caused the 500 errors?"`
)

// NewCommand creates a new cobra.Command for 'summarise' sub command
func NewCommand() *cobra.Command {
	command := v1summarise.Command{}

	cmd := &cobra.Command{
		Use:                   "summarise <environment> <stream> <stream>",
		Aliases:               []string{"summarize"},
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Summarise a window of an application's logs",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]

			if len(args) > 1 {
				command.Streams = args[1:]
			}

			return command.Run(cmd.Context())
		},
	}

	filter.AddFlags(cmd.Flags(), &command.Options)
	cmd.Flags().StringVar(&command.Prompt, "prompt", "", "Natural-language question for the summariser")

	return cmd
}
