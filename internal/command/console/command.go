package console

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"
	"github.com/skratchdot/open-golang/open"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config"
)

// Command for console access.
type Command struct {
	Environment string
	Print       bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	config, err := config.New()
	if err != nil {
		return err
	}

	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	cluster, err := client.Cluster().Get(ctx, &pb.ClusterGetRequest{})
	if err != nil {
		return fmt.Errorf("failed to get cluster information: %w", err)
	}
	if cluster.Endpoints == nil || cluster.Endpoints.Console == "" {
		return fmt.Errorf("cluster does not have a console domain")
	}
	consoleHost := cluster.Endpoints.Console

	_, err = client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	consoleURL := fmt.Sprintf("https://%s/projects/%s/%s/metrics", consoleHost, config.Project, cmd.Environment)

	if cmd.Print {
		fmt.Println(consoleURL)
		return nil
	}

	fmt.Println("Opening in Browser")

	return open.Run(consoleURL)
}
