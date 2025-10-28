package trace

import (
	"context"
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/api/pb"
	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/compass/cli/app"
	"github.com/skpr/compass/cli/app/events"
	applogger "github.com/skpr/compass/cli/app/logger"
	"github.com/skpr/compass/trace"
)

var (
	cmdLong = `
  Trace requests as they flow through your application using Compass.`

	cmdExample = `
  # Start profiling for the specified environment
  skpr trace <environment>`
)

// NewCommand creates a new cobra.Command for 'trace' sub command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "trace [environment]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Trace requests as they flow through your application",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			environment := args[0]

			p := tea.NewProgram(app.NewModel(""), tea.WithAltScreen())

			logger, err := applogger.New(p)
			if err != nil {
				return fmt.Errorf("failed to setup logger: %w", err)
			}

			ctx, cancel := context.WithCancel(cmd.Context())

			eg := errgroup.Group{}

			// Start the collector.
			eg.Go(func() error {
				logger.Info("Connecting to Skpr API...")

				client, _, err := wfclient.NewFromFile()
				if err != nil {
					return err
				}

				stream, err := client.Compass().StreamTraces(ctx, &pb.StreamTracesRequest{
					Environment: environment,
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
									StartTime: f.StartTime,
									Elapsed:   f.ElapsedTime,
								})
							}

							p.Send(events.Trace{
								IngestionTime: time.Unix(t.Metadata.StartTime, 0),
								Trace: trace.Trace{
									Metadata: trace.Metadata{
										RequestID: t.Metadata.RequestId,
										URI:       t.Metadata.Uri,
										Method:    t.Metadata.Method,
										StartTime: t.Metadata.StartTime,
										EndTime:   t.Metadata.EndTime,
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
		},
	}

	return cmd
}
