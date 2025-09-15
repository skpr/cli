package get

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command for getting a config.
type Command struct {
	Environment string
	Key         string
	ShowSecrets bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Config().Get(ctx, &pb.ConfigGetRequest{
		Name:       cmd.Environment,
		Key:        cmd.Key,
		ShowSecret: cmd.ShowSecrets,
	})
	if err != nil {
		return err
	}

	fmt.Print(resp.Config.Value)

	return nil
}
