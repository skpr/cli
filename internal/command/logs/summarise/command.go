package summarise

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/jwalton/gchalk"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/command/logs/filter"
)

// Command to summarise a window of an environment's logs.
type Command struct {
	filter.Options
	Prompt string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	logFilter, err := cmd.Build()
	if err != nil {
		return err
	}

	resp, err := client.Logs().Summarise(ctx, &pb.LogSummariseRequest{
		Filter: logFilter,
		Prompt: cmd.Prompt,
	})
	if err != nil {
		return fmt.Errorf("failed to summarise logs: %w", err)
	}

	print(os.Stdout, resp)

	return nil
}

// print renders the summary response to the given writer.
func print(w io.Writer, resp *pb.LogSummariseResponse) {
	heading := func(text string) string {
		return gchalk.WithHex(color.HexOrange).Bold(text)
	}

	if resp.Overview != "" {
		fmt.Fprintln(w, heading("Overview"))
		fmt.Fprintln(w, resp.Overview)
	}

	if len(resp.Bullets) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, heading("Notable events"))

		for _, bullet := range resp.Bullets {
			fmt.Fprintf(w, "  - %s\n", bullet)
		}
	}

	if len(resp.SuggestedActions) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, heading("Suggested actions"))

		for _, action := range resp.SuggestedActions {
			fmt.Fprintf(w, "  - %s\n", action)
		}
	}
}
