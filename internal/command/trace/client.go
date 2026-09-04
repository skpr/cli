package trace

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skpr/cli/internal/client"
)

type traceStream interface {
	Recv() (*pb.StreamTracesResponse, error)
}

type commandAPI interface {
	GetEnvironment(context.Context, *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error)
	StreamTraces(context.Context, *pb.StreamTracesRequest) (traceStream, error)
}

type connectFunc func(context.Context) (context.Context, commandAPI, error)

type skprAPI struct {
	client *client.Client
}

func connectAPI(ctx context.Context) (context.Context, commandAPI, error) {
	ctx, apiClient, err := client.New(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, skprAPI{client: apiClient}, nil
}

func (api skprAPI) GetEnvironment(ctx context.Context, request *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
	return api.client.Environment().Get(ctx, request)
}

func (api skprAPI) StreamTraces(ctx context.Context, request *pb.StreamTracesRequest) (traceStream, error) {
	return api.client.Trace().StreamTraces(ctx, request)
}

func (cmd *Command) preflight(ctx context.Context, connect connectFunc) (context.Context, commandAPI, error) {
	ctx, api, err := connect(ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to connect to Skpr API: %w", err)
	}

	_, err = api.GetEnvironment(ctx, &pb.EnvironmentGetRequest{Name: cmd.Environment})
	if status.Code(err) == codes.NotFound {
		return ctx, nil, fmt.Errorf("environment %q does not exist", cmd.Environment)
	}
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to verify environment %q: %w", cmd.Environment, err)
	}

	return ctx, api, nil
}
