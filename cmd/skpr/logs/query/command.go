package query

import (
	"time"

	"github.com/spf13/cobra"

	v1query "github.com/skpr/cli/internal/command/logs/query"
)

var (
	cmdLong = `Run a bounded query over an environment's logs.`

	cmdExample = `
  # Query the default streams of an environment over the last hour.
  skpr logs query dev

  # Query specific streams over the last 30 minutes.
  skpr logs query dev --stream nginx --stream fpm --since 30m

  # Query an absolute time range with substring filters.
  skpr logs query dev --from "1 hour ago" --to now --contains error --exclude healthcheck`
)

// NewCommand creates a new cobra.Command for 'query' sub command.
func NewCommand() *cobra.Command {
	command := v1query.Command{}

	cmd := &cobra.Command{
		Use:                   "query <environment>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Run a bounded query over an environment's logs",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Environment = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringSliceVar(&command.Streams, "stream", nil, "Stream to include in the query (repeatable)")
	cmd.Flags().DurationVar(&command.Since, "since", time.Hour, "Relative time window from now")
	cmd.Flags().StringVar(&command.From, "from", "", "Absolute start of the time range (used with --to)")
	cmd.Flags().StringVar(&command.To, "to", "", "Absolute end of the time range (used with --from)")
	cmd.Flags().StringSliceVar(&command.Contains, "contains", nil, "Substring an event must contain (repeatable)")
	cmd.Flags().StringSliceVar(&command.Exclude, "exclude", nil, "Substring an event must NOT contain (repeatable)")
	cmd.Flags().Int32Var(&command.Limit, "limit", 0, "Maximum number of events to return (0 for the server default)")
	cmd.Flags().BoolVar(&command.Indent, "indent", false, "Enable indenting for pretty printed logs")

	return cmd
}
