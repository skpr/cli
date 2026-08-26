package query

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/internal/command/logs/filter"
	v1query "github.com/skpr/cli/internal/command/logs/query"
)

var (
	cmdLong = `Run a bounded query over the logs of a running application.`

	cmdExample = `
  # Query the last hour of the default streams.
  skpr logs query dev

  # Query specific streams over the last 24 hours.
  skpr logs query dev nginx fpm --timeframe 24h

  # Query an absolute time range.
  skpr logs query dev --from "2 days ago" --to now

  # Query for events containing a substring, excluding another.
  skpr logs query dev --contains error --exclude healthcheck`
)

// NewCommand creates a new cobra.Command for 'query' sub command
func NewCommand() *cobra.Command {
	command := v1query.Command{}

	cmd := &cobra.Command{
		Use:                   "query <environment> [stream] [stream]",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Run a bounded query over the logs of a running application",
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
	cmd.Flags().Int32Var(&command.Limit, "limit", 0, "Cap on returned events (0 for server default)")
	cmd.Flags().BoolVar(&command.Indent, "indent", false, "Enable indenting for pretty printed logs")

	return cmd
}
