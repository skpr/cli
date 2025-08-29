package get

import (
	"fmt"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
)

// Command for getting a config.
type Command struct {
	Environment string
	Key         string
	ShowSecrets bool
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
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
