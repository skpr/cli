package summary

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	v1summary "github.com/skpr/cli/internal/command/logs/summary"
)

var (
	cmdLong = `Summarise an environment's logs using a natural-language prompt.`

	cmdExample = `
  # Summarise the last hour of logs with a question.
  skpr logs summary dev what errors happened recently

  # Quote the prompt when it contains shell metacharacters.
  skpr logs summary dev "why did fpm crash?"

  # Restrict to specific streams and a tighter window.
  skpr logs summary prod --stream nginx --since 30m top sources of 5xx responses`
)

// NewCommand creates a new cobra.Command for 'summary' sub command.
func NewCommand() *cobra.Command {
	command := v1summary.Command{}

	cmd := &cobra.Command{
		Use:                   "summary <environment> <prompt...>",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Summarise an environment's logs using a prompt",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			command.Prompt = strings.Join(args[1:], " ")
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringSliceVar(&command.Streams, "stream", nil, "Stream to include (repeatable)")
	cmd.Flags().DurationVar(&command.Since, "since", time.Hour, "Relative time window from now")
	cmd.Flags().StringVar(&command.From, "from", "", "Absolute start of the time range (used with --to)")
	cmd.Flags().StringVar(&command.To, "to", "", "Absolute end of the time range (used with --from)")
	cmd.Flags().StringSliceVar(&command.Contains, "contains", nil, "Substring an event must contain (repeatable)")
	cmd.Flags().StringSliceVar(&command.Exclude, "exclude", nil, "Substring an event must NOT contain (repeatable)")

	return cmd
}
