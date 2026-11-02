package rsync

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1rsync "github.com/skpr/cli/internal/command/rsync"
)

var (
	cmdLong = `Sync files from remote and local directories.

Private files (/mnt/private) may contain personal information and are not downloaded by default.
Syncing from a parent directory, such as /mnt/, excludes them automatically. Syncing from the
private filesystem itself fails unless --include-private is set. Only download private files to
investigate a specific issue, with approval.`

	cmdExample = `
  # Sync from local to a remote environment
  skpr rsync ./foo/ dev:/mnt/temporary/foo/

  # Download public files from a remote environment
  skpr rsync dev:/mnt/public/ ./files/

  # Download private files from a remote environment (asks for confirmation)
  skpr rsync --include-private dev:/mnt/private/ ./private/`
)

// NewCommand creates a new cobra.Command for 'rsync' sub command
func NewCommand() *cobra.Command {
	command := v1rsync.Command{}

	cmd := &cobra.Command{
		Use:                   "rsync <source> <destination>",
		DisableFlagsInUseLine: true,
		Short:                 "Sync files between local and remote environments",
		Args:                  cobra.ExactArgs(2),
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupSecureShell,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Source = args[0]
			command.Destination = args[1]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringSliceVar(&command.Excludes, "exclude", []string{}, "Exclude files matching PATTERN.")
	cmd.Flags().StringVar(&command.ExcludeFrom, "exclude-from", "", "Exclude files matching patterns in FILE.")
	cmd.Flags().BoolVar(&command.DryRun, "dry-run", false, "Skpr will not request a confirmation for configuration changes.")
	cmd.Flags().BoolVar(&command.IncludePrivate, "include-private", false, "Download private files ("+v1rsync.PrivateFilesPath+"). Asks for confirmation.")
	cmd.Flags().BoolVarP(&command.Force, "force", "f", false, "Skip the confirmation prompt when downloading private files.")

	return cmd
}
