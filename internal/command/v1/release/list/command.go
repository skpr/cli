package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// Command that lists all releases.
type Command struct {
	Params struct {
		JSON bool
	}
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	releases, err := client.Release().List(ctx, &pb.ReleaseListRequest{})
	if err != nil {
		return errors.Wrap(err, "Could not list releases")
	}

	if len(releases.Items) == 0 {
		fmt.Println("No releases found:", "See `skpr package` to create a new release")
		return nil
	}

	header := []string{
		"Name",
		"Date",
		"Environments",
	}

	var rows [][]string

	for _, item := range releases.Items {
		rows = append(rows, []string{
			item.Name,
			item.Date,
			strings.Join(item.Environments, "\n"),
		})
	}

	return table.Print(os.Stdout, header, rows)
}
