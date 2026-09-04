package trace

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/pkg/app/events"
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
			retryDelay = traceRetryInitial
			err = receiveTraces(ctx, stream, sender)
		}

		if ctx.Err() != nil {
			return nil
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

func receiveTraces(ctx context.Context, stream traceStream, sender messageSender) error {
	for {
		response, err := stream.Recv()
		if err != nil {
			return err
		}
		if response == nil {
			return errors.New("received an empty trace stream response")
		}

		for _, item := range response.GetTraces() {
			converted := traceFromProto(item)
			sender.Send(events.Trace{
				IngestionTime: converted.Metadata.StartTime,
				Trace:         converted,
			})
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
