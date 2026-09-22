package suspend

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
	Suspend(context.Context, *pb.TraceSuspendRequest, ...grpc.CallOption) (*pb.TraceSuspendResponse, error)
}

type connectFunc func(context.Context) (context.Context, api, error)

func connectAPI(ctx context.Context) (context.Context, api, error) {
	ctx, apiClient, err := client.New(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, apiClient.Trace(), nil
}

// Command which suspends tracing for an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	return cmd.run(ctx, connectAPI, os.Stderr)
}

func (cmd *Command) run(ctx context.Context, connect connectFunc, w io.Writer) error {
	ctx, api, err := connect(ctx)
	if err != nil {
		return err
	}

	_, err = api.Suspend(ctx, &pb.TraceSuspendRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to suspend tracing for environment %q: %w", cmd.Environment, err)
	}

	fmt.Fprintf(w, "Tracing has been suspended for environment %q.\n", cmd.Environment)

	return nil
}
