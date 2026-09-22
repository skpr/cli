package suspend

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type fakeAPI struct {
	suspend func(*pb.TraceSuspendRequest) (*pb.TraceSuspendResponse, error)
}

func (api *fakeAPI) Suspend(_ context.Context, request *pb.TraceSuspendRequest, _ ...grpc.CallOption) (*pb.TraceSuspendResponse, error) {
	return api.suspend(request)
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

func TestSuspend(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		suspend: func(request *pb.TraceSuspendRequest) (*pb.TraceSuspendResponse, error) {
			assert.Equal(t, "dev", request.GetEnvironment())
			return &pb.TraceSuspendResponse{}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Equal(t, "Tracing has been suspended for environment \"dev\".\n", b.String())
}

func TestSuspendFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		suspend: func(*pb.TraceSuspendRequest) (*pb.TraceSuspendResponse, error) {
			return nil, errors.New("permission denied")
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to suspend tracing for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestSuspendConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), func(ctx context.Context) (context.Context, api, error) {
		return ctx, nil, errors.New("connection refused")
	}, &b)

	require.EqualError(t, err, "connection refused")
}
