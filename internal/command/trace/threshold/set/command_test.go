package set

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
)

type fakeAPI struct {
	setThreshold func(*pb.TraceSetThresholdRequest) (*pb.TraceSetThresholdResponse, error)
}

func (api *fakeAPI) SetThreshold(_ context.Context, request *pb.TraceSetThresholdRequest, _ ...grpc.CallOption) (*pb.TraceSetThresholdResponse, error) {
	return api.setThreshold(request)
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

func TestSetThreshold(t *testing.T) {
	cmd := Command{Environment: "dev", Threshold: "10ms"}
	api := &fakeAPI{
		setThreshold: func(request *pb.TraceSetThresholdRequest) (*pb.TraceSetThresholdResponse, error) {
			assert.Equal(t, "dev", request.GetEnvironment())
			assert.Equal(t, 10*time.Millisecond, request.GetThreshold().AsDuration())
			return &pb.TraceSetThresholdResponse{}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Equal(t, "Tracing threshold for environment \"dev\" has been set to 10ms.\n", b.String())
}

func TestSetThresholdZero(t *testing.T) {
	cmd := Command{Environment: "dev", Threshold: "0s"}
	api := &fakeAPI{
		setThreshold: func(request *pb.TraceSetThresholdRequest) (*pb.TraceSetThresholdResponse, error) {
			assert.Equal(t, time.Duration(0), request.GetThreshold().AsDuration())
			return &pb.TraceSetThresholdResponse{}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
}

func TestSetThresholdInvalid(t *testing.T) {
	for _, threshold := range []string{"", "soon", "10", "-10ms"} {
		t.Run(threshold, func(t *testing.T) {
			cmd := Command{Environment: "dev", Threshold: threshold}

			var b bytes.Buffer

			// Connecting would panic, which is how we know a malformed duration
			// is rejected before the API is reached.
			err := cmd.run(context.Background(), func(context.Context) (context.Context, api, error) {
				panic("connected despite a malformed threshold")
			}, &b)

			require.Error(t, err)
			assert.Contains(t, err.Error(), threshold)
			assert.Empty(t, b.String())
		})
	}
}

func TestSetThresholdFailure(t *testing.T) {
	cmd := Command{Environment: "dev", Threshold: "10ms"}
	api := &fakeAPI{
		setThreshold: func(*pb.TraceSetThresholdRequest) (*pb.TraceSetThresholdResponse, error) {
			return nil, errors.New("permission denied")
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to set tracing threshold for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestSetThresholdConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "dev", Threshold: "10ms"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), func(ctx context.Context) (context.Context, api, error) {
		return ctx, nil, errors.New("connection refused")
	}, &b)

	require.EqualError(t, err, "connection refused")
}
