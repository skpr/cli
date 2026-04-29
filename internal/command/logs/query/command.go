package query

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	logshared "github.com/skpr/cli/internal/command/logs"
)

// Format used for printing event timestamps.
const timestampFormat = "15:04:05.000"

// Separator drawn between columns.
var separator = gchalk.Dim("│")

// Command runs a bounded log query.
type Command struct {
	Environment string
	Streams     []string
	Since       time.Duration
	From        string
	To          string
	Contains    []string
	Exclude     []string
	Limit       int32
	Indent      bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	filter, err := logshared.BuildFilter(logshared.FilterParams{
		Environment: cmd.Environment,
		Streams:     cmd.Streams,
		Since:       cmd.Since,
		From:        cmd.From,
		To:          cmd.To,
		Contains:    cmd.Contains,
		Exclude:     cmd.Exclude,
	})
	if err != nil {
		return err
	}

	ctx, c, err := client.New(ctx)
	if err != nil {
		return err
	}

	stream, err := c.Logs().Query(ctx, &pb.LogQueryRequest{
		Filter: filter,
		Limit:  cmd.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to start query: %w", err)
	}

	var meta *pb.LogQueryMeta

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive query response: %w", err)
		}

		switch body := resp.GetBody().(type) {
		case *pb.LogQueryResponse_Batch:
			for _, ev := range body.Batch.GetEvents() {
				printEvent(ev, cmd.Indent)
			}
		case *pb.LogQueryResponse_Meta:
			meta = body.Meta
		}
	}

	if meta != nil {
		ranAt := meta.GetRanAt().AsTime().Format(time.RFC3339)
		fmt.Printf("\nScanned %d events at %s\n", meta.GetScanned(), ranAt)
	}

	return nil
}

func printEvent(ev *pb.LogEvent, indent bool) {
	ts := gchalk.Dim(ev.GetTimestamp().AsTime().Format(timestampFormat))

	streamLabel := formatStream(ev.GetStream())

	message := logshared.PrettyPrint(ev.GetMessage(), indent)

	fmt.Printf("%s %s %s %s %s\n", ts, separator, streamLabel, separator, message)
}

// formatStream colourises the visible stream name and pads/truncates to width.
// Trailing pad spaces are kept outside the colour codes to avoid background bleed.
func formatStream(name string) string {
	return gchalk.WithHex(color.HexOrange).Bold(name)
}
