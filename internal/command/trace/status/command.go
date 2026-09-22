package status

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// thresholdNotReported is shown when an environment does not report a tracing
// threshold. A zero threshold traces every call, so it cannot stand in for one
// which was never reported.
const thresholdNotReported = "Not reported"

// api is the subset of the trace service which this command uses.
type api interface {
	GetSuspended(context.Context, *pb.TraceGetSuspendedRequest, ...grpc.CallOption) (*pb.TraceGetSuspendedResponse, error)
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

// Command which reports the tracing status for an environment.
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

	suspended, err := api.GetSuspended(ctx, &pb.TraceGetSuspendedRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to get tracing suspended state for environment %q: %w", cmd.Environment, err)
	}

	threshold, err := api.GetThreshold(ctx, &pb.TraceGetThresholdRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to get tracing threshold for environment %q: %w", cmd.Environment, err)
	}

	return Print(w, cmd.Environment, suspended.GetSuspended(), threshold.GetThreshold().AsDuration(), threshold.GetThreshold() != nil)
}

// Print the tracing status for an environment.
//
// The threshold is only shown when the environment reported one, which
// hasThreshold declares. Unlike "skpr trace threshold get", a threshold which
// was never reported does not fail the command, so that the suspended state is
// still shown.
func Print(w io.Writer, environment string, suspended bool, threshold time.Duration, hasThreshold bool) error {
	header := []string{
		"Property",
		"Value",
	}

	tracing := "Active"
	if suspended {
		tracing = "Suspended"
	}

	formattedThreshold := thresholdNotReported
	if hasThreshold {
		formattedThreshold = threshold.String()
	}

	rows := [][]string{
		{"Environment", environment},
		{"Tracing", tracing},
		{"Threshold", formattedThreshold},
	}

	return table.Print(w, header, rows)
}
