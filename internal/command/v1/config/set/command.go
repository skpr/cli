package set

import (
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command for setting config.
type Command struct {
	Environment string
	Key         string
	Value       string
	FromFile    string
	Secret      bool
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
		_, err := client.Config().Get(ctx, &pb.ConfigGetRequest{
			Name: cmd.Environment,
			Key:  cmd.Key,
		})
		if err == nil {
			if ok := confirmation.Confirm(cmd.Force, "Key already exists, are you sure you want to replace it? [yes/no]"); !ok {
				return nil
			}
		}
	}

	req := &pb.ConfigSetRequest{
		Name: cmd.Environment,
		Config: &pb.Config{
			Key:    cmd.Key,
			Value:  cmd.Value,
			Secret: cmd.Secret,
		},
	}

	if cmd.FromFile != "" {
		data, err := os.ReadFile(cmd.FromFile)
		if err != nil {
			return fmt.Errorf("failed to load value from file: %w", err)
		}

		req.Config.Value = string(data)
	}

	if req.Config.Value == "" {
		return fmt.Errorf("value was not provided")
	}

	_, err = client.Config().Set(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
