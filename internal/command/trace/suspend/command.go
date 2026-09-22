package suspend

import (
	"context"
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command which suspends tracing for an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	_, err = client.Trace().Suspend(ctx, &pb.TraceSuspendRequest{Environment: cmd.Environment})
	if err != nil {
		return fmt.Errorf("failed to suspend tracing for environment %q: %w", cmd.Environment, err)
	}

	fmt.Fprintf(os.Stderr, "Tracing has been suspended for environment %q.\n", cmd.Environment)

	return nil
}
