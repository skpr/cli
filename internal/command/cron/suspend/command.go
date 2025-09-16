package suspend

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command to suspend cron jobs.
type Command struct {
	Environment string
	Wait        bool
	Timeout     time.Duration
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	_, err = client.Cron().Suspend(ctx, &pb.CronSuspendRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if cmd.Wait {
		fmt.Fprintln(os.Stderr, "Waiting for cron to finish")

		limiter := time.Tick(5 * time.Second)
		finalFinishTime := time.Now().Add(cmd.Timeout)

		for {
			<-limiter

			if time.Now().After(finalFinishTime) {
				fmt.Fprintln(os.Stderr, "Review the long-running jobs and restart.\nSuspension remains active.")
				break
			}

			resp, err := client.Cron().JobList(ctx, &pb.CronJobListRequest{
				Environment: cmd.Environment,
			})
			if err != nil {
				return fmt.Errorf("failed to list jobs: %w", err)
			}

			status := new(bool)
			for _, cron := range resp.List {
				if cron.Phase == pb.CronJobDetail_Running {
					*status = true
					fmt.Fprintf(os.Stderr, "Still waiting for '%v' to finish (%v)\n", cron.Name, cron.Duration)
				}
			}

			if !*status {
				break
			}
		}
	}

	fmt.Fprintln(os.Stderr, "All Cron tasks have been suspended.")

	return nil
}
