package query

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/TylerBrock/colorjson"
	faithcolor "github.com/fatih/color"
	"github.com/jwalton/gchalk"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/command/logs/filter"
)

// Command to run a bounded query over an environment's logs.
type Command struct {
	filter.Options
	Limit  int32
	Indent bool
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

	stream, err := client.Logs().Query(ctx, &pb.LogQueryRequest{
		Filter: logFilter,
		Limit:  cmd.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to run query: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to receive query response: %w", err)
		}

		if batch := resp.GetBatch(); batch != nil {
			for _, event := range batch.Events {
				printEvent(event, cmd.Indent)
			}
		}

		if meta := resp.GetMeta(); meta != nil {
			fmt.Fprintf(os.Stderr, "Scanned %d events\n", meta.Scanned)
		}
	}

	return nil
}

// printEvent prints a single log event prefixed with its timestamp and stream.
func printEvent(event *pb.LogEvent, indent bool) {
	prefix := gchalk.WithHex(color.HexOrange).Bold(event.Stream)
	timestamp := event.Timestamp.AsTime().Format(time.RFC3339)
	message := prettyPrint(event.Message, indent)

	fmt.Println(timestamp, prefix, message)
}

// Returns a pretty output for JSON messages.
func prettyPrint(message string, indent bool) string {
	var obj map[string]interface{}

	err := json.Unmarshal([]byte(message), &obj)
	if err != nil {
		return message
	}

	formatter := colorjson.NewFormatter()
	formatter.KeyColor = faithcolor.New(faithcolor.FgWhite).Add(faithcolor.Bold)

	if indent {
		formatter.Indent = 2
	}

	raw, err := formatter.Marshal(obj)
	if err != nil {
		return message
	}

	return string(raw)
}
