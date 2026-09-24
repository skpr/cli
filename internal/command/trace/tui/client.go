package tui

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
	GetSuspended(context.Context, *pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error)
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

func (api skprAPI) GetSuspended(ctx context.Context, request *pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
	return api.client.Trace().GetSuspended(ctx, request)
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

	// The trace stream is rejected for an environment which is not collecting
	// traces, so check before starting the app rather than letting it retry.
	suspended, err := api.GetSuspended(ctx, &pb.TraceGetSuspendedRequest{Environment: cmd.Environment})
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to verify tracing for environment %q: %s", cmd.Environment, status.Convert(err).Message())
	}
	if suspended.GetSuspended() {
		return ctx, nil, fmt.Errorf("tracing is suspended for environment %q, run \"skpr trace resume %s\" to resume it", cmd.Environment, cmd.Environment)
	}

	return ctx, api, nil
}
