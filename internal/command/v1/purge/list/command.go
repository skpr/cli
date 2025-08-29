package list

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/list"
)

// Command to list purge requests.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Purge().List(ctx, &pb.PurgeListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return errors.Wrap(err, "Could not list purge requests")
	}

	table := uitable.New()
	table.Wrap = true

	table.AddRow("ID", "CREATED", "PATH", "STATUS")

	for _, request := range resp.Requests {
		paths, err := list.Print(request.Paths)
		if err != nil {
			return errors.Wrap(err, "Unable to render paths list")
		}

		table.AddRow(request.ID, request.Created, paths, request.Status)
	}

	fmt.Println(table.String())

	return nil
}
