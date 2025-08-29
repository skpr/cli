package dashboard

import (
	"fmt"

	"github.com/skratchdot/open-golang/open"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
)

// Command for dashboard access.
type Command struct {
	Environment string
	Print       bool
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	respGet, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	if respGet.Environment.Dashboard == nil || respGet.Environment.Dashboard.URL == "" {
		return fmt.Errorf("environment does not have a dashboard")
	}

	if cmd.Print {
		fmt.Println(respGet.Environment.Dashboard.URL)
		return nil
	}

	fmt.Println("Opening in Browser")

	return open.Run(respGet.Environment.Dashboard.URL)
}
