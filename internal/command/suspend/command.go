package suspend

import (
	"context"
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command to suspend cron jobs.
type Command struct {
	Environment string
	Force       bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	env, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if env.Environment.Production {
		if ok := confirmation.Confirm(cmd.Force, "Are you sure you want to suspend a production environment? [yes/no]"); !ok {
			return nil
		}
	}

	_, err = client.Environment().Suspend(ctx, &pb.EnvironmentSuspendRequest{
		Name: cmd.Environment,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "The environment has been suspended.")

	return nil
}
