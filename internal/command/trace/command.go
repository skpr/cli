package trace

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/tracing/cli/app"
	"github.com/skpr/compass/tracing/cli/app/events"
	applogger "github.com/skpr/compass/tracing/cli/app/logger"
	"github.com/skpr/compass/tracing/trace"
	"golang.org/x/sync/errgroup"
	"io"

	"github.com/skpr/cli/internal/client"
)

// Command to trace environments.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	p := tea.NewProgram(app.NewModel(""), tea.WithAltScreen())

	logger, err := applogger.New(p)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	eg := errgroup.Group{}

	// Start the collector.
	eg.Go(func() error {
		logger.Info("Connecting to Skpr API...")

		ctx, client, err := client.New(ctx)
		if err != nil {
			return err
		}

		stream, err := client.Trace().StreamTraces(ctx, &pb.StreamTracesRequest{
			Environment: cmd.Environment,
		})
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("streaming trace failed: %w", err)
				}

				for _, t := range resp.Traces {
					var fcalls []trace.FunctionCall

					for _, f := range t.FunctionCalls {
						fcalls = append(fcalls, trace.FunctionCall{
							Name:      f.Name,
							StartTime: int64(f.StartTime.Nanos),
							Elapsed:   int64(f.Elapsed.Nanos),
						})
					}

					p.Send(events.Trace{
						IngestionTime: t.Metadata.StartTime.AsTime(),
						Trace: trace.Trace{
							Metadata: trace.Metadata{
								RequestID: t.Metadata.RequestId,
								URI:       t.Metadata.Uri,
								Method:    t.Metadata.Method,
								StartTime: int64(t.Metadata.StartTime.Nanos),
								EndTime:   t.Metadata.EndTime.Seconds,
							},
							FunctionCalls: fcalls,
						},
					})
				}
			}
		}
	})

	// Start the application.
	eg.Go(func() error {
		_, err := p.Run()
		if err != nil {
			return fmt.Errorf("failed to run program: %w", err)
		}

		cancel()

		return nil
	})

	return eg.Wait()
}
