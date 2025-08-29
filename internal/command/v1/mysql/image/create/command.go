package create

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
)

// Command to create an image.
type Command struct {
	Database    string
	Environment string
	Tag         string
	Policies    []string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	resp, err := client.Mysql().ImageCreate(ctx, &pb.ImageCreateRequest{
		Environment: cmd.Environment,
		Database:    cmd.Database,
		Tag:         cmd.Tag,
		Policies:    cmd.Policies,
	})
	if err != nil {
		return errors.Wrap(err, "image creation failed")
	}

	for _, image := range resp.Images {
		fmt.Println("Building new mysql image:", image.ID)
	}

	return nil
}
