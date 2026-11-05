package resume

import (
	"context"
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command for resuming cron jobs.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	_, err = client.Environment().Resume(ctx, &pb.EnvironmentResumeRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "The environment has been resumed.")

	return nil
}
