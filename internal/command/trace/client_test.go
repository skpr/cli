package trace

import (
	"context"
	"errors"
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeCommandAPI struct {
	getEnvironment func(context.Context, *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error)
	streamTraces   func(context.Context, *pb.StreamTracesRequest) (traceStream, error)
}

func (api *fakeCommandAPI) GetEnvironment(ctx context.Context, request *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
	return api.getEnvironment(ctx, request)
}

func (api *fakeCommandAPI) StreamTraces(ctx context.Context, request *pb.StreamTracesRequest) (traceStream, error) {
	return api.streamTraces(ctx, request)
}

func TestPreflightConnectionFailure(t *testing.T) {
	cmd := Command{Environment: "staging"}

	_, api, err := cmd.preflight(context.Background(), func(ctx context.Context) (context.Context, commandAPI, error) {
		return ctx, nil, errors.New("connection refused")
	})

	require.EqualError(t, err, "failed to connect to Skpr API: connection refused")
	assert.Nil(t, api)
}

func TestPreflightEnvironmentNotFound(t *testing.T) {
	cmd := Command{Environment: "missing"}
	api := &fakeCommandAPI{
		getEnvironment: func(_ context.Context, request *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
			assert.Equal(t, "missing", request.GetName())
			return nil, status.Error(codes.NotFound, "environment not found")
		},
	}

	_, validatedAPI, err := cmd.preflight(context.Background(), func(ctx context.Context) (context.Context, commandAPI, error) {
		return ctx, api, nil
	})

	require.EqualError(t, err, `environment "missing" does not exist`)
	assert.Nil(t, validatedAPI)
}

func TestPreflightEnvironmentLookupFailure(t *testing.T) {
	cmd := Command{Environment: "staging"}
	api := &fakeCommandAPI{
		getEnvironment: func(context.Context, *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		},
	}

	_, validatedAPI, err := cmd.preflight(context.Background(), func(ctx context.Context) (context.Context, commandAPI, error) {
		return ctx, api, nil
	})

	require.EqualError(t, err, `failed to verify environment "staging": rpc error: code = PermissionDenied desc = permission denied`)
	assert.Nil(t, validatedAPI)
}

func TestPreflightSuccess(t *testing.T) {
	type contextKey string
	const key contextKey = "authenticated"

	cmd := Command{Environment: "staging"}
	api := &fakeCommandAPI{
		getEnvironment: func(ctx context.Context, request *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
			assert.Equal(t, true, ctx.Value(key))
			assert.Equal(t, "staging", request.GetName())
			return &pb.EnvironmentGetResponse{}, nil
		},
	}

	ctx, validatedAPI, err := cmd.preflight(context.Background(), func(ctx context.Context) (context.Context, commandAPI, error) {
		return context.WithValue(ctx, key, true), api, nil
	})

	require.NoError(t, err)
	assert.Equal(t, true, ctx.Value(key))
	assert.Same(t, api, validatedAPI)
}
