package list

import (
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/components/tooltip"
	skprtable "github.com/skpr/cli/internal/table"
)

// Helpful text provided by the tooltip.
const tooltipText = `To view logs for a specific stream, use the tail command.

$ skpr logs tail ENVIRONMENT STREAM STREAM`

// Command to list all log sources.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Logs().ListStreams(ctx, &pb.LogListStreamsRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to list streams: %w", err)
	}

	header := []string{
		"Streams",
	}

	var rows [][]string

	for _, stream := range resp.Streams {
		rows = append(rows, []string{stream})
	}

	err = skprtable.Print(os.Stdout, header, rows)
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	err = tooltip.Render(os.Stdout, tooltipText)
	if err != nil {
		return fmt.Errorf("failed to render tooltip: %w", err)
	}

	return nil
}

// Row which can be....
type Row struct {
	Stream string `header:"stream"`
}
