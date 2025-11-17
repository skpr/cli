package set

import (
	"context"
	"fmt"
	"strings"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command for setting config.
type Command struct {
	Key   string
	Value string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	switch cmd.Key {
	case "contact":
		req := &pb.SetContactRequest{
			Contact: cmd.Value,
		}
		_, err = client.Project().SetContact(ctx, req)
		if err != nil {
			return err
		}
	case "tags":
		req := &pb.SetTagsRequest{
			Tags: strings.Fields(cmd.Value),
		}
		_, err = client.Project().SetTags(ctx, req)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid key. Currently supported keys are: contact, tags")
	}

	return nil
}
