package console

import (
	"context"
	"fmt"
	"strings"

	"github.com/skpr/api/pb"
	"github.com/skratchdot/open-golang/open"

	"github.com/skpr/cli/internal/client"
	clientconfig "github.com/skpr/cli/internal/client/config"
)

// Command for console access.
type Command struct {
	Environment string
	Print       bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	config, err := clientconfig.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build the console URL by replacing the "cluster." prefix with "console."
	host := config.API.Host()
	consoleHost := strings.Replace(host, "cluster.", "console.", 1)
	// @todo Remove the /metrics when we respond correctly in UI.
	consoleURL := fmt.Sprintf("https://%s/projects/%s/%s/metrics", consoleHost, config.Project, cmd.Environment)

	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	_, err = client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	if cmd.Print {
		fmt.Println(consoleURL)
		return nil
	}

	fmt.Println("Opening in Browser")

	return open.Run(consoleURL)
}
