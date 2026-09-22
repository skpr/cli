package resume

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
	resume func(*pb.TraceResumeRequest) (*pb.TraceResumeResponse, error)
}

func (api *fakeAPI) Resume(_ context.Context, request *pb.TraceResumeRequest, _ ...grpc.CallOption) (*pb.TraceResumeResponse, error) {
	return api.resume(request)
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

func TestResume(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		resume: func(request *pb.TraceResumeRequest) (*pb.TraceResumeResponse, error) {
			assert.Equal(t, "dev", request.GetEnvironment())
			return &pb.TraceResumeResponse{}, nil
		},
	}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(api), &b))
	assert.Equal(t, "Tracing has been resumed for environment \"dev\".\n", b.String())
}

func TestResumeFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}
	api := &fakeAPI{
		resume: func(*pb.TraceResumeRequest) (*pb.TraceResumeResponse, error) {
			return nil, errors.New("permission denied")
		},
	}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(api), &b)
	require.EqualError(t, err, `failed to resume tracing for environment "dev": permission denied`)
	assert.Empty(t, b.String())
}

func TestResumeConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), func(ctx context.Context) (context.Context, api, error) {
		return ctx, nil, errors.New("connection refused")
	}, &b)

	require.EqualError(t, err, "connection refused")
}
