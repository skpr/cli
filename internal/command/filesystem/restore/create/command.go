package create

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command to create a volume restore task.
type Command struct {
	Environment string
	Backup      string
	VolumeName  string
	Wait        bool
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
		if ok := confirmation.Confirm(cmd.Force, "Are you sure you want to restore a backup to production? [yes/no]"); !ok {
			return nil
		}
	}

	fmt.Fprintln(os.Stderr, "Creating new restore")

	resp, err := client.Volume().RestoreCreate(ctx, &pb.VolumeRestoreCreateRequest{
		Environment: cmd.Environment,
		Backup:      cmd.Backup,
		VolumeName:  cmd.VolumeName,
	})
	if err != nil {
		return fmt.Errorf("failed to create restore: %w", err)
	}

	fmt.Println(resp.ID)

	if !cmd.Wait {
		return nil
	}

	fmt.Fprintln(os.Stderr, "Waiting for restore to finish")

	limiter := time.Tick(10 * time.Second)

	for {
		<-limiter

		resp, err := client.Volume().RestoreGet(ctx, &pb.VolumeRestoreGetRequest{
			ID: resp.ID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get volume restore")
		}

		switch resp.VolumeRestore.Phase {
		case pb.VolumeRestoreStatus_Completed:
			fmt.Fprintln(os.Stderr, "Restore complete!")
			return nil
		case pb.VolumeRestoreStatus_Failed:
			return fmt.Errorf("the restore failed: the Skpr team has been notified")
		case pb.VolumeRestoreStatus_Unknown:
			return fmt.Errorf("the restore failed for an unknown reason: the Skpr team has been notified")
		}
	}
}
