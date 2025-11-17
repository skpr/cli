package get

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/project"
	"github.com/skpr/cli/internal/table"
)

// Command for getting a config.
type Command struct {
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Project().Get(ctx, &pb.ProjectGetRequest{})
	if err != nil {
		return err
	}

	err = Print(os.Stdout, resp.Project)
	if err != nil {
		return err
	}

	return nil
}

// Print the table...
func Print(w io.Writer, item *pb.Project) error {
	header := []string{
		"Attribute",
		"Value",
	}

	rows := [][]string{
		{"ID", item.ID},
		{"Name", item.Name},
		{"Contact", item.Contact},
		{"Version", item.Version},
		{"Environments", strings.Join(project.ListEnvironmentsByName(item), ", ")},
		{"Size", item.Size},
		{"Tags", strings.Join(item.Tags, ", ")},
	}

	return table.Print(w, header, rows)
}
