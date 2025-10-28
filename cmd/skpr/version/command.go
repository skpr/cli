package version

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1version "github.com/skpr/cli/internal/command/version"
)

var (
	cmdLong = `Print client and server version information`

	// GitVersion overridden at build time by:
	//   -ldflags="-X github.com/skpr/cli/internal/command/v1/version.GitVersion=${VERSION}"
	GitVersion string
	// BuildDate overridden at build time by:
	//   -ldflags="-X github.com/skpr/cli/internal/command/v1/version.BuildDate=${BUILD_DATE}"
	// BuildDate = time.Now().Format("2006-01-02")
	BuildDate string
)

// Options is the commandline options for 'version' sub command
type Options struct {
	Debug bool
}

// NewCommand creates a new cobra.Command for 'version' sub command
func NewCommand() *cobra.Command {
	command := v1version.Command{}

	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Print client and server version information",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupOther,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.Run(cmd.Context(), GitVersion, BuildDate)
		},
	}

	cmd.Flags().BoolVarP(&command.Debug, "Debug", "d", false, "Turn on debugging when interacting with the server.")

	return cmd
}
