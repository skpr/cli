package get

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
	getThreshold func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error)
}

func (api *fakeAPI) GetThreshold(_ context.Context, request *pb.TraceGetThresholdRequest, _ ...grpc.CallOption) (*pb.TraceGetThresholdResponse, error) {
	return api.getThreshold(request)
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

func TestGetThreshold(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getThreshold: func(request *pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			assert.Equal(t, "dev", request.GetEnvironment())
			return &pb.TraceGetThresholdResponse{Threshold: durationpb.New(10 * time.Millisecond)}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Equal(t, "10ms\n", b.String())
}

// A zero threshold traces every call, so it is a meaningful value rather than a
// stand in for one the environment never reported.
func TestGetThresholdZero(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return &pb.TraceGetThresholdResponse{Threshold: durationpb.New(0)}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Equal(t, "0s\n", b.String())
}

func TestGetThresholdMissing(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return &pb.TraceGetThresholdResponse{}, nil
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `environment "dev" did not report a tracing threshold`)
	assert.Empty(t, b.String())
}

func TestGetThresholdFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		getThreshold: func(*pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
			return nil, errors.New("permission denied")
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to get tracing threshold for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestGetThresholdConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), func(ctx context.Context) (context.Context, api, error) {
		return ctx, nil, errors.New("connection refused")
	}, &b)

	require.EqualError(t, err, "connection refused")
}
