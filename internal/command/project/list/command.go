package set

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/project"
	"github.com/skpr/cli/internal/table"
)

// Command for listing projects.
type Command struct {
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Project().List(ctx, &pb.ProjectListRequest{})
	if err != nil {
		return errors.Wrap(err, "Could not list projects")
	}

	return Print(os.Stdout, resp.Projects)
}

// Print the table...
func Print(w io.Writer, list []*pb.Project) error {
	if len(list) == 0 {
		fmt.Fprintln(w, "No projects found")
		return nil
	}

	sortProjects(list)

	header := []string{
		"Name",
		"Contact",
		"Version",
		"Environments",
		"Size",
		"Tags",
	}

	var rows [][]string

	for _, item := range list {
		rows = append(rows, []string{
			item.Name,
			item.Contact,
			item.Version,
			strings.Join(project.ListEnvironmentsByName(item), ", "),
			item.Size,
			strings.Join(item.Tags, ", "),
		})
	}

	return table.Print(w, header, rows)
}

// sortProjects sorts the list of projects by name in ascending order.
func sortProjects(list []*pb.Project) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
}
