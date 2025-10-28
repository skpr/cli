package create

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// Command to create a filesystem backup.
type Command struct {
	Environment string
	VolumeName  string
	Wait        bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Creating new filesystem backup")

	resp, err := client.Volume().BackupCreate(ctx, &pb.VolumeBackupCreateRequest{
		Environment: cmd.Environment,
		VolumeName:  cmd.VolumeName,
	})
	if err != nil {
		return fmt.Errorf("failed to create filesystem backup: %w", err)
	}

	fmt.Println(resp.ID)

	if !cmd.Wait {
		return nil
	}

	fmt.Fprintln(os.Stderr, "Waiting for filesystem backup to finish")

	limiter := time.Tick(10 * time.Second)

	for {
		<-limiter

		resp, err := client.Volume().BackupGet(ctx, &pb.VolumeBackupGetRequest{
			ID: resp.ID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get filesystem backup")
		}

		switch resp.VolumeBackup.Phase {
		case pb.VolumeBackupStatus_Completed:
			fmt.Fprintln(os.Stderr, "Backup complete!")
			return nil
		case pb.VolumeBackupStatus_Failed:
			return fmt.Errorf("the filesystem backup failed: the Skpr team has been notified")
		case pb.VolumeBackupStatus_Unknown:
			return fmt.Errorf("the filesystem backup failed for an unknown reason: the Skpr team has been notified")
		}
	}
}
