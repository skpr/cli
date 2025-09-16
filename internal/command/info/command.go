package info

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command to get information about the environment.
type Command struct {
	Name   string
	Format string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{Name: cmd.Name})
	if err != nil {
		return errors.Wrap(err, "Could not fetch environment")
	}

	switch cmd.Format {
	case "json":
		out, err := json.MarshalIndent(resp.Environment, "", "\t")
		if err != nil {
			return errors.Wrap(err, "Could not marshal environment to json")
		}
		fmt.Println(string(out))
	default:
		return fmt.Errorf("invalid output format specified - \"%s\"", cmd.Format)
	}

	return nil
}
