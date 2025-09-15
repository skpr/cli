package suspend

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command that suspends daemons.
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

	_, err = client.Daemon().Suspend(ctx, &pb.DaemonSuspendRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if cmd.Wait {
		fmt.Fprintln(os.Stderr, "Waiting for daemon to finish")

		limiter := time.Tick(5 * time.Second)
		finalFinishTime := time.Now().Add(cmd.Timeout)

		for {
			<-limiter

			if time.Now().After(finalFinishTime) {
				fmt.Fprintln(os.Stderr, "Review the long-running daemons and restart.\nSuspension remains active.")
				break
			}

			resp, err := client.Daemon().List(ctx, &pb.DaemonListRequest{
				Environment: cmd.Environment,
			})
			if err != nil {
				return fmt.Errorf("failed to list daemons: %w", err)
			}

			var status bool
			for _, daemon := range resp.List {
				if !daemon.Suspended {
					status = true
					fmt.Fprintf(os.Stderr, "Still waiting for '%s' to finish\n", daemon.Name)
				}
			}

			if !status {
				break
			}
		}
	}

	fmt.Fprintln(os.Stderr, "All Daemon tasks have been suspended.")

	return nil
}
