package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// Command for listing cronjobs.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Cron().List(ctx, &pb.CronListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if len(resp.List) == 0 {
		fmt.Fprintln(os.Stderr, "No CronJobs were found.")
	}

	header := []string{
		"Name",
		"Schedule",
		"Command",
		"Last Schedule",
		"Last Successful Execution",
		"Suspended",
	}

	var rows [][]string

	for _, item := range resp.List {
		rows = append(rows, []string{
			item.Name,
			item.Schedule,
			item.Command,
			item.LastScheduleTime,
			item.LastSuccessfulTime,
			strconv.FormatBool(item.Suspended),
		})
	}

	return table.Print(os.Stdout, header, rows)
}
