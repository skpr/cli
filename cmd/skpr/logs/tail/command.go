package tail

import (
	"github.com/spf13/cobra"

	v1tail "github.com/skpr/cli/internal/command/logs/tail"
)

var (
	cmdLong = `Tail the logs of a running application.`

	cmdExample = `
  # Tail the logs of a running application using the default stream.
  skpr logs tail dev
  
  # Tail the logs of a running application using multiple streams.
  skpr logs tail dev nginx fpm`
)

// NewCommand creates a new cobra.Command for 'tail' sub command
func NewCommand() *cobra.Command {
	command := v1tail.Command{}

	cmd := &cobra.Command{
		Use:                   "tail <environment> <stream> <stream>",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Tail the logs of a running application",
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

	cmd.Flags().BoolVar(&command.Indent, "indent", false, "Enable indenting for pretty printed logs")

	return cmd
}
