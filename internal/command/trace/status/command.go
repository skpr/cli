package status

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/table"
)

// thresholdNotReported is shown when an environment does not report a tracing
// threshold. A zero threshold traces every call, so it cannot stand in for one
// which was never reported.
const thresholdNotReported = "Not reported"

// Command which reports the tracing status for an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	suspended, err := client.Trace().GetSuspended(ctx, &pb.TraceGetSuspendedRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to get tracing suspended state for environment %q: %w", cmd.Environment, err)
	}

	threshold, err := client.Trace().GetThreshold(ctx, &pb.TraceGetThresholdRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to get tracing threshold for environment %q: %w", cmd.Environment, err)
	}

	var duration *time.Duration

	if threshold.GetThreshold() != nil {
		value := threshold.GetThreshold().AsDuration()
		duration = &value
	}

	return Print(os.Stdout, cmd.Environment, suspended.GetSuspended(), duration)
}

// Print the tracing status for an environment.
//
// A nil threshold is one the environment did not report. Unlike
// "skpr trace threshold get", that does not fail the command, so that the
// suspended state is still shown.
func Print(w io.Writer, environment string, suspended bool, threshold *time.Duration) error {
	header := []string{
		"Property",
		"Value",
	}

	tracing := "Active"
	if suspended {
		tracing = "Suspended"
	}

	formattedThreshold := thresholdNotReported
	if threshold != nil {
		formattedThreshold = threshold.String()
	}

	rows := [][]string{
		{"Environment", environment},
		{"Tracing", tracing},
		{"Threshold", formattedThreshold},
	}

	return table.Print(w, header, rows)
}
