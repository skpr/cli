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

// Command to create a backup.
type Command struct {
	Environment  string
	DatabaseName string
	Wait         bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Creating new mysql backup")

	resp, err := client.Mysql().BackupCreate(ctx, &pb.MysqlBackupCreateRequest{
		Environment:  cmd.Environment,
		DatabaseName: cmd.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create mysql backup: %w", err)
	}

	fmt.Println(resp.ID)

	if !cmd.Wait {
		return nil
	}

	fmt.Fprintln(os.Stderr, "Waiting for mysql backup to finish")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		resp, err := client.Backup().Get(ctx, &pb.BackupGetRequest{
			ID: resp.ID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get backup")
		}

		switch resp.Backup.Phase {
		case pb.BackupStatus_Completed:
			fmt.Fprintln(os.Stderr, "Backup complete!")
			return nil
		case pb.BackupStatus_Failed:
			return fmt.Errorf("the mysql backup failed: the Skpr team has been notified")
		case pb.BackupStatus_Unknown:
			return fmt.Errorf("the mysql backup failed for an unknown reason: the Skpr team has been notified")
		}
	}
}
