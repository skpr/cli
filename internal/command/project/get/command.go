package get

import (
	"context"
	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/table"
	"io"
	"os"
	"strings"

	"github.com/skpr/cli/internal/client"
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
func Print(w io.Writer, project *pb.Project) error {
	header := []string{
		"Attribute",
		"Value",
	}

	rows := [][]string{
		{"ID", project.ID},
		{"Name", project.Name},
		{"Version", project.Version},
		{"Tags", strings.Join(project.Tags, ", ")},
		{"Contact", project.Contact},
	}

	return table.Print(w, header, rows)
}
