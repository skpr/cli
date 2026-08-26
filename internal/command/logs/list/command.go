package list

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/command/logs/streams"
	"github.com/skpr/cli/internal/components/tooltip"
	skprtable "github.com/skpr/cli/internal/table"
)

// Helpful text provided by the tooltip.
const tooltipText = `To tail a stream in real-time, use the tail command.

$ skpr logs tail ENVIRONMENT STREAM STREAM

To run a bounded query over a time window, use the query command.

$ skpr logs query ENVIRONMENT STREAM STREAM`

// Command to list all log sources.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := streams.List(ctx, client.Logs(), cmd.Environment)
	if err != nil {
		return fmt.Errorf("failed to list streams: %w", err)
	}

	header := []string{
		"Stream",
		"Capabilities",
	}

	var rows [][]string

	for _, stream := range resp.Streams {
		var capabilities []string

		for _, streamType := range stream.Types {
			capabilities = append(capabilities, streamType.String())
		}

		rows = append(rows, []string{stream.Name, strings.Join(capabilities, ", ")})
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
