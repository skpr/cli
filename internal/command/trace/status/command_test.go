package status

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

type fakeAPI struct {
	getSuspended func(*pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error)
	getThreshold func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error)
}

func (api *fakeAPI) GetSuspended(_ context.Context, request *pb.TraceGetSuspendedRequest, _ ...grpc.CallOption) (*pb.TraceGetSuspendedResponse, error) {
	return api.getSuspended(request)
}

func (api *fakeAPI) GetThreshold(_ context.Context, request *pb.TraceGetThresholdRequest, _ ...grpc.CallOption) (*pb.TraceGetThresholdResponse, error) {
	return api.getThreshold(request)
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

func TestStatus(t *testing.T) {
	for _, tt := range []struct {
		name      string
		suspended bool
		expected  string
	}{
		{name: "active", suspended: false, expected: "Active"},
		{name: "suspended", suspended: true, expected: "Suspended"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command{Environment: "dev"}
			api := &fakeAPI{
				getSuspended: func(request *pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
					assert.Equal(t, "dev", request.GetEnvironment())
					return &pb.TraceGetSuspendedResponse{Suspended: tt.suspended}, nil
				},
				getThreshold: func(request *pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
					assert.Equal(t, "dev", request.GetEnvironment())
					return &pb.TraceGetThresholdResponse{Threshold: durationpb.New(10 * time.Millisecond)}, nil
				},
			}

			var b bytes.Buffer

			require.NoError(t, cmd.run(context.Background(), connect(api), &b))
			assert.Contains(t, b.String(), "dev")
			assert.Contains(t, b.String(), tt.expected)
			assert.Contains(t, b.String(), "10ms")
		})
	}
}

// A zero threshold traces every call, so it is a meaningful value rather than a
// stand in for one the environment never reported.
func TestStatusThresholdZero(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getSuspended: func(*pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
			return &pb.TraceGetSuspendedResponse{}, nil
		},
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return &pb.TraceGetThresholdResponse{Threshold: durationpb.New(0)}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Contains(t, b.String(), "0s")
	assert.NotContains(t, b.String(), thresholdNotReported)
}

// Unlike "skpr trace threshold get", a threshold which was never reported is
// shown rather than failing the command, so that the suspended state is still
// reported.
func TestStatusThresholdMissing(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getSuspended: func(*pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
			return &pb.TraceGetSuspendedResponse{Suspended: true}, nil
		},
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return &pb.TraceGetThresholdResponse{}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Contains(t, b.String(), thresholdNotReported)
	assert.Contains(t, b.String(), "Suspended")
}

func TestStatusSuspendedFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getSuspended: func(*pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
			return nil, errors.New("permission denied")
		},
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			t.Fatal("threshold was read after the suspended state failed")
			return nil, nil
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to get tracing suspended state for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestStatusThresholdFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getSuspended: func(*pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
			return &pb.TraceGetSuspendedResponse{}, nil
		},
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return nil, errors.New("permission denied")
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to get tracing threshold for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestStatusConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), func(ctx context.Context) (context.Context, api, error) {
		return ctx, nil, errors.New("connection refused")
	}, &b)

	require.EqualError(t, err, "connection refused")
}
