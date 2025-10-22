package list

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
	timeutils "github.com/skpr/cli/internal/time"
)

// Command to list images.
type Command struct {
	Environment string
}

// Row used for formatting list response.
type Row struct {
	ID             string
	Phase          string
	StartTime      string
	CompletionTime string
	Duration       string
	Tags           []string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	resp, err := client.Mysql().ImageList(ctx, &pb.ImageListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return errors.Wrap(err, "image list failed")
	}

	if len(resp.List) == 0 {
		fmt.Println("No images found:", "See `skpr mysql image create` to create an image")
		return nil
	}

	var rows []Row

	for _, item := range resp.List {
		row := Row{
			ID:    item.ID,
			Phase: item.Phase.String(),
			Tags:  item.Tags,
		}

		if item.StartTime != "" {
			start, err := timeutils.ParseString(item.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time: %w", err)
			}

			row.StartTime = start.Format(time.RFC822Z)
		}

		if item.CompletionTime != "" {
			completion, err := timeutils.ParseString(item.CompletionTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time: %w", err)
			}

			row.CompletionTime = completion.Format(time.RFC822Z)

			row.Duration = item.Duration
		}

		rows = append(rows, row)
	}

	header := []string{
		"ID",
		"Phase",
		"Start Time",
		"Completion Time",
		"Duration",
		"Tags",
	}

	var flatRows [][]string

	for _, item := range rows {
		flatRows = append(flatRows, []string{
			item.ID,
			color.ApplyColorToString(item.Phase),
			item.StartTime,
			item.CompletionTime,
			item.Duration,
			strings.Join(item.Tags, "\n"),
		})
	}

	return table.Print(os.Stdout, header, flatRows)
}
