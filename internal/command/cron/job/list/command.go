package list

import (
	"context"
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
)

// Command for listing jobs.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Cron().JobList(ctx, &pb.CronJobListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if len(resp.List) == 0 {
		fmt.Fprintln(os.Stderr, "No Cron Jobs were found.")
	}

	header := []string{
		"Name",
		"Phase",
		"Start Time",
		"Duration",
	}

	var rows [][]string

	for _, item := range resp.List {
		rows = append(rows, []string{
			item.Name,
			color.ApplyColorToString(item.Phase.String()),
			item.StartTime,
			item.Duration,
		})
	}

	return table.Print(os.Stdout, header, rows)
}
