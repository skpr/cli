package set

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/skpr/cli/internal/client"
)

// api is the subset of the trace service which this command uses.
type api interface {
	SetThreshold(context.Context, *pb.TraceSetThresholdRequest, ...grpc.CallOption) (*pb.TraceSetThresholdResponse, error)
}

type connectFunc func(context.Context) (context.Context, api, error)

func connectAPI(ctx context.Context) (context.Context, api, error) {
	ctx, apiClient, err := client.New(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, apiClient.Trace(), nil
}

// Command which sets the tracing threshold for an environment.
//
// Threshold is the value to apply, as a duration which time.ParseDuration
// accepts.
type Command struct {
	Environment string
	Threshold   string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	return cmd.run(ctx, connectAPI, os.Stderr)
}

func (cmd *Command) run(ctx context.Context, connect connectFunc, w io.Writer) error {
	// Parse before connecting so that a malformed duration fails immediately
	// rather than after a round trip to the API.
	threshold, err := time.ParseDuration(cmd.Threshold)
	if err != nil {
		return fmt.Errorf("failed to parse threshold %q: %w", cmd.Threshold, err)
	}

	if threshold < 0 {
		return fmt.Errorf("threshold %q cannot be negative", cmd.Threshold)
	}

	ctx, api, err := connect(ctx)
	if err != nil {
		return err
	}

	_, err = api.SetThreshold(ctx, &pb.TraceSetThresholdRequest{
		Environment: cmd.Environment,
		Threshold:   durationpb.New(threshold),
	})
	if err != nil {
		return fmt.Errorf("failed to set tracing threshold for environment %q: %w", cmd.Environment, err)
	}

	fmt.Fprintf(w, "Tracing threshold for environment %q has been set to %s.\n", cmd.Environment, threshold)

	return nil
}
