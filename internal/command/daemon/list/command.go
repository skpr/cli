package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// Command that lists daemons.
type Command struct {
	Environment string
	JSON        bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Daemon().List(ctx, &pb.DaemonListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if len(resp.List) == 0 {
		fmt.Fprintln(os.Stderr, "No DaemonJobs were found.")
	}

	if cmd.JSON {
		data, err := json.Marshal(resp.List)

		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return Print(os.Stdout, resp.List)
}

// Print the table...
func Print(w io.Writer, list []*pb.DaemonDetail) error {
	header := []string{
		"Name",
		"Command",
		"Suspended",
	}

	var rows [][]string

	for _, item := range list {
		rows = append(rows, []string{
			item.Name,
			item.Command,
			strconv.FormatBool(item.Suspended),
		})
	}

	return table.Print(w, header, rows)
}
