package get

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc"

	"github.com/skpr/cli/internal/client"
)

// api is the subset of the trace service which this command uses.
type api interface {
	GetThreshold(context.Context, *pb.TraceGetThresholdRequest, ...grpc.CallOption) (*pb.TraceGetThresholdResponse, error)
}

type connectFunc func(context.Context) (context.Context, api, error)

func connectAPI(ctx context.Context) (context.Context, api, error) {
	ctx, apiClient, err := client.New(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, apiClient.Trace(), nil
}

// Command which gets the tracing threshold for an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	return cmd.run(ctx, connectAPI, os.Stdout)
}

func (cmd *Command) run(ctx context.Context, connect connectFunc, w io.Writer) error {
	ctx, api, err := connect(ctx)
	if err != nil {
		return err
	}

	resp, err := api.GetThreshold(ctx, &pb.TraceGetThresholdRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to get tracing threshold for environment %q: %w", cmd.Environment, err)
	}

	// A zero threshold traces every call, so it cannot stand in for one which
	// was never reported.
	if resp.GetThreshold() == nil {
		return fmt.Errorf("environment %q did not report a tracing threshold", cmd.Environment)
	}

	fmt.Fprintln(w, resp.GetThreshold().AsDuration())

	return nil
}
