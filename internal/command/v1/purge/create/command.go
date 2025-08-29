package create

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
)

// Command to create a purge request.
type Command struct {
	Environment string
	Paths       []string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Purge().Create(ctx, &pb.PurgeCreateRequest{
		Environment: cmd.Environment,
		Paths:       cmd.Paths,
	})
	if err != nil {
		return errors.Wrap(err, "Could not create purge request")
	}

	fmt.Println("Invalidation submitted:", resp.ID)

	return nil
}
