package logs

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/logs/list"
	"github.com/skpr/cli/cmd/skpr/logs/query"
	"github.com/skpr/cli/cmd/skpr/logs/summarise"
	"github.com/skpr/cli/cmd/skpr/logs/tail"
	skprcommand "github.com/skpr/cli/internal/command"
)

var (
	cmdLong = `Debug an application using logs.`
)

// NewCommand creates a new cobra.Command for 'logss' sub command
func NewCommand(featureSummary bool) *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "logs",
		DisableFlagsInUseLine: true,
		Short:                 "Debug an application using logs.",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupDebug,
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(query.NewCommand())
	cmd.AddCommand(tail.NewCommand())

	if featureSummary {
		cmd.AddCommand(summarise.NewCommand())
	}

	return cmd
}
