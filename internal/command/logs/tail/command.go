package tail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/TylerBrock/colorjson"
	faithcolor "github.com/fatih/color"
	"github.com/jwalton/gchalk"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
)

// Command to stream the logs for an environment.
type Command struct {
	Environment string
	Streams     []string
	Indent      bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	list, err := client.Logs().ListStreams(ctx, &pb.LogListStreamsRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to get stream list: %w", err)
	}

	// If none are provided use the default stream provided by the API.
	if len(cmd.Streams) == 0 && list.Default != "" {
		cmd.Streams = []string{
			list.Default,
		}
	}

	if len(cmd.Streams) == 0 {
		return fmt.Errorf("no streams provided and no default stream found")
	}

	// Validate that the provided streams are correct.
	for _, stream := range cmd.Streams {
		if !slices.Contains(list.Streams, stream) {
			return fmt.Errorf("stream not found: %s", stream)
		}
	}

	fmt.Println("Following streams:", strings.Join(cmd.Streams, ", "))

	e := errgroup.Group{}

	for _, stream := range cmd.Streams {
		stream := stream

		prefix := gchalk.WithHex(color.HexOrange).Bold(stream)

		e.Go(func() error {
			tail, err := client.Logs().Tail(ctx, &pb.LogTailRequest{
				Environment: cmd.Environment,
				Stream:      stream,
			})
			if err != nil {
				return fmt.Errorf("failed to initiate stream: %s: %w", stream, err)
			}

			for {
				resp, err := tail.Recv()
				if err == io.EOF {
					break
				}

				if err != nil {
					return fmt.Errorf("fail to tail stream: %s: %w", stream, err)
				}

				message := prettyPrint(resp.Message, cmd.Indent)

				// Only prefix when there is more than one stream.
				if len(cmd.Streams) > 1 {
					fmt.Println(prefix, message)
				} else {
					fmt.Println(message)
				}
			}

			return nil
		})
	}

	return e.Wait()
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
