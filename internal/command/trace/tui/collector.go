package tui

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/pkg/app/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	traceRetryInitial = time.Second
	traceRetryMaximum = 30 * time.Second
)

type messageSender interface {
	Send(tea.Msg)
}

type errorLogger interface {
	Error(string, ...any)
}

func collectTraces(ctx context.Context, api commandAPI, environment string, sender messageSender, logger errorLogger) error {
	retryDelay := traceRetryInitial

	for {
		if ctx.Err() != nil {
			return nil
		}

		sender.Send(events.Connection{State: events.ConnectionStateConnecting})

		stream, err := api.StreamTraces(ctx, &pb.StreamTracesRequest{Environment: environment})
		if err == nil && stream == nil {
			err = errors.New("trace stream was not created")
		}
		if err == nil {
			sender.Send(events.Connection{State: events.ConnectionStateConnected})

			// Opening a stream does not wait on the server, so a rejected stream
			// only fails on its first receive. The backoff is only reset once
			// traces arrive, otherwise every attempt would retry after the
			// initial delay.
			var received bool
			received, err = receiveTraces(ctx, stream, sender)
			if received {
				retryDelay = traceRetryInitial
			}
		}

		if ctx.Err() != nil {
			return nil
		}

		// The server rejects the stream when the environment is not collecting
		// traces, eg. tracing was suspended, which retrying will not fix.
		if status.Code(err) == codes.FailedPrecondition {
			return fmt.Errorf("trace stream rejected: %s", status.Convert(err).Message())
		}

		if errors.Is(err, io.EOF) {
			err = errors.New("trace stream closed")
		} else {
			err = fmt.Errorf("trace stream failed: %w", err)
		}

		logger.Error(err.Error())
		sender.Send(events.Connection{State: events.ConnectionStateRetrying, Err: err})

		timer := time.NewTimer(retryDelay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return nil
		case <-timer.C:
		}

		if retryDelay < traceRetryMaximum {
			retryDelay *= 2
			if retryDelay > traceRetryMaximum {
				retryDelay = traceRetryMaximum
			}
		}
	}
}

func receiveTraces(ctx context.Context, stream traceStream, sender messageSender) (bool, error) {
	var received bool

	for {
		response, err := stream.Recv()
		if err != nil {
			return received, err
		}
		if response == nil {
			return received, errors.New("received an empty trace stream response")
		}

		received = true

		for _, item := range response.GetTraces() {
			converted := traceFromProto(item)
			sender.Send(events.Trace{
				IngestionTime: converted.Metadata.StartTime,
				Trace:         converted,
			})
		}

		if ctx.Err() != nil {
			return received, ctx.Err()
		}
	}
}
