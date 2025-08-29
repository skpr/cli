package delete

import (
	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command for deleting a config.
type Command struct {
	Environment string
	Key         string
	Force       bool
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	env, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Environment,
	})

	if err != nil {
		return err
	}

	if env.Environment.Production {
		if ok := confirmation.Confirm(cmd.Force, "Are you sure you want to PERMANENTLY DELETE this production config? [yes/no]"); !ok {
			return nil
		}
	}

	_, err = client.Config().Delete(ctx, &pb.ConfigDeleteRequest{
		Name: cmd.Environment,
		Key:  cmd.Key,
	})
	if err != nil {
		return err
	}

	return nil
}
