package list

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
)

// Command to list environments.
type Command struct{}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Environment().List(ctx, &pb.EnvironmentListRequest{})
	if err != nil {
		return errors.Wrap(err, "Could not list environments")
	}

	return Print(os.Stdout, resp.Environments)
}

// Print the table...
func Print(w io.Writer, list []*pb.Environment) error {
	if len(list) == 0 {
		fmt.Fprintln(w, "No environments found:", "See `skpr create` to provision a new environment")
		return nil
	}

	sortEnvs(list)

	header := []string{
		"Name",
		"Version",
		"Size",
		"Routes",
		"Phase",
	}

	var rows [][]string

	for _, item := range list {
		rows = append(rows, []string{
			item.Name,
			item.Version,
			item.Size,
			strings.Join(append(item.Ingress.Routes, item.Ingress.Domain), "\n"),
			color.ApplyColorToString(item.Phase),
		})
	}

	return table.Print(w, header, rows)
}

// sortEnvs sorts the list putting the production envs last.
func sortEnvs(list []*pb.Environment) {
	// Ensure prod environments are listed last.
	sort.Slice(list, func(i, j int) bool {
		if list[i].Production != list[j].Production {
			return !list[i].Production
		}
		return list[i].Name < list[j].Name
	})
}
