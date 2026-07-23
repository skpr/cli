package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/info"
	"github.com/skpr/cli/cmd/skpr/list"
	"github.com/skpr/cli/cmd/skpr/login"
	"github.com/skpr/cli/cmd/skpr/logs"
	"github.com/skpr/cli/cmd/skpr/mysql"
	"github.com/skpr/cli/cmd/skpr/version"
	"github.com/skpr/cli/internal/client/config/user"
	skprcommand "github.com/skpr/cli/internal/command"
)

const cmdExample = `
	# List all environments.
	skpr list

    # Show information about an environment.
    skpr info ENVIRONMENT_NAME

	# List all log streams for an environment. Used to find the name of the log stream you want to stream.
	skpr logs list ENVIRONMENT_NAME

	# Stream logs for an environment. Replace LOG_STREAM_NAME with the name of the log stream you want to stream.
	# Requires a Ctrl+C to stop streaming.
	skpr logs stream ENVIRONMENT_NAME LOG_STREAM_NAME LOG_STREAM_NAME

	# Pull MySQL database from an environment to your local machine.
	skpr mysql pull ENVIRONMENT_NAME
`

var cmd = &cobra.Command{
	Use:     "skpr",
	Short:   "Command line interface for interacting with the Skpr Hosting Platform.",
	Example: cmdExample,
}

func main() {
	// Load our configuration which contains aliases and feature flags.
	userConfig, err := user.NewClient()
	if err != nil {
		fmt.Println("Failed to load user config file:", err)
		os.Exit(1)
	}

	// Experimental commands.
	featureFlags, err := userConfig.LoadFeatureFlags()
	if err != nil {
		fmt.Println("Failed to load feature flags:", err)
		os.Exit(1)
	}

	skprcommand.AddGroupsToCommand(cmd)

	cmd.AddCommand(info.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(login.NewCommand())
	cmd.AddCommand(logs.NewCommand())
	cmd.AddCommand(mysql.NewCommand(featureFlags.DockerClient))
	cmd.AddCommand(version.NewCommand())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
