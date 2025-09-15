package info

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// Command to print release info.
type Command struct {
	JSON    bool
	Version string
}

const (
	// ManifestTypeRuntime is the type of runtime image.
	ManifestTypeRuntime = "runtime"
)

type manifestItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Release().Info(ctx, &pb.ReleaseInfoRequest{
		Name: cmd.Version,
	})
	if err != nil {
		return fmt.Errorf("could not get release: %w", err)
	}

	if cmd.JSON {
		var data []manifestItem

		for _, image := range resp.Images {
			data = append(data, manifestItem{
				Name: image.Name,
				Type: ManifestTypeRuntime,
				Tag:  image.URI,
			})
		}

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(b))
		return nil
	}

	return Print(os.Stdout, resp)
}

// Print the table...
func Print(w io.Writer, item *pb.ReleaseInfoResponse) error {
	header := []string{
		"Service",
		"Date",
		"Image",
	}

	var rows [][]string

	for _, image := range item.Images {
		rows = append(rows, []string{
			image.Name,
			item.Date,
			image.URI,
		})
	}

	return table.Print(os.Stdout, header, rows)
}
