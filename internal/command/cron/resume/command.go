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

	_, err = client.Cron().Resume(ctx, &pb.CronResumeRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "All Cron tasks have been resumed.")

	return nil
}
