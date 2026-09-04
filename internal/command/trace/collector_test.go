package trace

import (
	"context"
	"errors"
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/pkg/app/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTraceStream struct {
	recv func() (*pb.StreamTracesResponse, error)
}

func (stream *fakeTraceStream) Recv() (*pb.StreamTracesResponse, error) {
	return stream.recv()
}

type recordingSender struct {
	messages []tea.Msg
	onSend   func(tea.Msg)
}

func (sender *recordingSender) Send(message tea.Msg) {
	sender.messages = append(sender.messages, message)
	if sender.onSend != nil {
		sender.onSend(message)
	}
}

type recordingLogger struct {
	errors []string
}

func (logger *recordingLogger) Error(message string, _ ...any) {
	logger.errors = append(logger.errors, message)
}

func TestCollectTracesSurfacesStreamErrors(t *testing.T) {
	tests := []struct {
		name           string
		streamTraces   func(context.Context, *pb.StreamTracesRequest) (traceStream, error)
		expectedError  string
		expectedStates []events.ConnectionState
	}{
		{
			name: "stream setup",
			streamTraces: func(_ context.Context, request *pb.StreamTracesRequest) (traceStream, error) {
				assert.Equal(t, "staging", request.GetEnvironment())
				return nil, errors.New("unavailable")
			},
			expectedError: "trace stream failed: unavailable",
			expectedStates: []events.ConnectionState{
				events.ConnectionStateConnecting,
				events.ConnectionStateRetrying,
			},
		},
		{
			name: "stream receive",
			streamTraces: func(context.Context, *pb.StreamTracesRequest) (traceStream, error) {
				return &fakeTraceStream{recv: func() (*pb.StreamTracesResponse, error) {
					return nil, errors.New("connection lost")
				}}, nil
			},
			expectedError: "trace stream failed: connection lost",
			expectedStates: []events.ConnectionState{
				events.ConnectionStateConnecting,
				events.ConnectionStateConnected,
				events.ConnectionStateRetrying,
			},
		},
		{
			name: "stream EOF",
			streamTraces: func(context.Context, *pb.StreamTracesRequest) (traceStream, error) {
				return &fakeTraceStream{recv: func() (*pb.StreamTracesResponse, error) {
					return nil, io.EOF
				}}, nil
			},
			expectedError: "trace stream closed",
			expectedStates: []events.ConnectionState{
				events.ConnectionStateConnecting,
				events.ConnectionStateConnected,
				events.ConnectionStateRetrying,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			api := &fakeCommandAPI{streamTraces: tt.streamTraces}
			logger := &recordingLogger{}
			sender := &recordingSender{}
			sender.onSend = func(message tea.Msg) {
				if connection, ok := message.(events.Connection); ok && connection.State == events.ConnectionStateRetrying {
					cancel()
				}
			}

			err := collectTraces(ctx, api, "staging", sender, logger)
			require.NoError(t, err)
			assert.Equal(t, []string{tt.expectedError}, logger.errors)

			var states []events.ConnectionState
			for _, message := range sender.messages {
				if connection, ok := message.(events.Connection); ok {
					states = append(states, connection.State)
					if connection.State == events.ConnectionStateRetrying {
						require.EqualError(t, connection.Err, tt.expectedError)
					}
				}
			}
			assert.Equal(t, tt.expectedStates, states)
		})
	}
}

func TestCollectTracesCancellationIsClean(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sender := &recordingSender{}
	logger := &recordingLogger{}
	err := collectTraces(ctx, &fakeCommandAPI{}, "staging", sender, logger)

	require.NoError(t, err)
	assert.Empty(t, sender.messages)
	assert.Empty(t, logger.errors)
}

func TestCollectTracesActiveStreamCancellationIsClean(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	api := &fakeCommandAPI{
		streamTraces: func(context.Context, *pb.StreamTracesRequest) (traceStream, error) {
			return &fakeTraceStream{recv: func() (*pb.StreamTracesResponse, error) {
				cancel()
				return nil, context.Canceled
			}}, nil
		},
	}

	sender := &recordingSender{}
	logger := &recordingLogger{}
	err := collectTraces(ctx, api, "staging", sender, logger)

	require.NoError(t, err)
	assert.Empty(t, logger.errors)
	assert.Equal(t, []tea.Msg{
		events.Connection{State: events.ConnectionStateConnecting},
		events.Connection{State: events.ConnectionStateConnected},
	}, sender.messages)
}
