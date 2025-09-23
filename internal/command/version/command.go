package version

import (
	"context"
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/version"
)

// Command that print the client and server versions.
type Command struct {
	Debug bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context, gitVersion, buildDate string) error {
	params := version.PrintParams{
		ClientVersion:   gitVersion,
		ClientBuildDate: buildDate,
	}

	// Get server version if we are in a project directory.
	ctx, client, err := client.New(ctx)
	if err == nil {
		resp, err := client.Version().Get(ctx, &pb.VersionGetRequest{})
		if err != nil && cmd.Debug {
			return err
		}

		if resp != nil {
			params.ServerVersion = resp.Version
			params.ServerBuildDate = resp.BuildDate
		}
	}

	if cmd.Debug && err != nil {
		fmt.Println(err)
	}

	return version.Print(os.Stdin, params)
}
