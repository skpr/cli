package create

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	skprlog "github.com/skpr/cli/internal/log"
	"github.com/skpr/cli/internal/retry"
)

// Command for creating backups.
type Command struct {
	Environment    string
	Wait           bool
	WaitTimeout    time.Duration
	WaitErrorLimit int32
	MySQL          MySQL
}

// MySQL parameter list provided by developer.
type MySQL struct {
	Policies []string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	prettyHandler := skprlog.NewHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})

	logger := slog.New(prettyHandler)

	logger.Info("Creating a new backup")

	resp, err := client.Backup().Create(ctx, &pb.BackupCreateRequest{
		Environment: cmd.Environment,
		MySQL: &pb.BackupCreateRequestMySQL{
			Policies: cmd.MySQL.Policies,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Println(resp.ID)

	if !cmd.Wait {
		return nil
	}

	logger.Info("Waiting for backup to finish")

	retryClient, err := retry.New(
		retry.WithLogger(logger),
		retry.WithTimeout(cmd.WaitTimeout),
		retry.WithErrorLimit(cmd.WaitErrorLimit),
	)
	if err != nil {
		return err
	}

	return retryClient.Poll(ctx, wait(ctx, logger, client.Backup(), cmd.Environment, resp.ID))
}

// Wait for the restore to complete.
func wait(ctx context.Context, logger *slog.Logger, client pb.BackupClient, environmentName, backupName string) retry.PollFunc {
	return func() (bool, error) {
		listResp, err := client.List(ctx, &pb.BackupListRequest{
			Environment: environmentName,
		})
		if err != nil {
			return false, fmt.Errorf("failed to get backup: %w", err)
		}

		backup := getBackup(backupName, listResp.List)

		if backup == nil {
			return true, fmt.Errorf("backup does not exist")
		}

		switch backup.Phase {
		case pb.BackupStatus_Completed:
			logger.Info("Backup complete!")
			return true, nil
		case pb.BackupStatus_Failed:
			return true, fmt.Errorf("the backup failed: the Skpr team has been notified")
		case pb.BackupStatus_Unknown:
			return true, fmt.Errorf("the backup failed for an unknown reason: the Skpr team has been notified")
		}

		return false, nil
	}
}

// Returns a restore status if one exists in the list.
func getBackup(name string, list []*pb.BackupStatus) *pb.BackupStatus {
	for _, item := range list {
		if item.Name == name {
			return item
		}
	}

	return nil
}
