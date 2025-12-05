package events

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
)

// Command to delete the environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	list, err := client.Events().List(ctx, &pb.EventsListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	return Print(os.Stdout, list.Events)
}

// Print the table...
func Print(w io.Writer, list []*pb.Event) error {
	if len(list) == 0 {
		fmt.Fprintln(w, "No events found")
		return nil
	}

	// Sort by date descending
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Timestamp.AsTime().After(list[j].Timestamp.AsTime())
	})

	header := []string{
		"Date",
		"Severity",
		"Type",
		"Message",
	}

	var rows [][]string

	for _, item := range list {
		rows = append(rows, []string{
			item.Timestamp.AsTime().Local().Format(time.RFC1123),
			color.ApplyColorToString(item.Severity.String()),
			item.Type,
			item.Message,
		})
	}

	return table.Print(w, header, rows)
}
