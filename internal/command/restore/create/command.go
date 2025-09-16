package create

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
	skprlog "github.com/skpr/cli/internal/log"
	"github.com/skpr/cli/internal/retry"
)

// Command to create a restore.
type Command struct {
	Environment    string
	Backup         string
	Wait           bool
	WaitTimeout    time.Duration
	WaitErrorLimit int32
	Force          bool
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

	prettyHandler := skprlog.NewHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})

	logger := slog.New(prettyHandler)

	logger.Info("Creating new restore")

	resp, err := client.Restore().Create(ctx, &pb.RestoreCreateRequest{
		Environment: cmd.Environment,
		Backup:      cmd.Backup,
	})
	if err != nil {
		return fmt.Errorf("failed to create restore: %w", err)
	}

	fmt.Println(resp.ID)

	if !cmd.Wait {
		return nil
	}

	logger.Info("Waiting for restore to finish")

	retryClient, err := retry.New(
		retry.WithLogger(logger),
		retry.WithTimeout(cmd.WaitTimeout),
		retry.WithErrorLimit(cmd.WaitErrorLimit),
	)
	if err != nil {
		return err
	}

	return retryClient.Poll(ctx, wait(ctx, logger, client.Restore(), cmd.Environment, resp.ID))
}

// Wait for the restore to complete.
func wait(ctx context.Context, logger *slog.Logger, client pb.RestoreClient, environmentName, restoreName string) retry.PollFunc {
	return func() (bool, error) {
		listResp, err := client.List(ctx, &pb.RestoreListRequest{
			Environment: environmentName,
		})
		if err != nil {
			return false, fmt.Errorf("failed to get backup: %w", err)
		}

		restore := getRestore(restoreName, listResp.List)

		if restore == nil {
			return true, fmt.Errorf("backup does not exist")
		}

		switch restore.Phase {
		case pb.RestoreStatus_Completed:
			logger.Info("Restore complete!")
			return true, nil
		case pb.RestoreStatus_Failed:
			return true, fmt.Errorf("the restore failed: the Skpr team has been notified")
		case pb.RestoreStatus_Unknown:
			return true, fmt.Errorf("the restore failed for an unknown reason: the Skpr team has been notified")
		}

		return false, nil
	}
}

// Returns a restore status if one exists in the list.
func getRestore(name string, list []*pb.RestoreStatus) *pb.RestoreStatus {
	for _, item := range list {
		if item.Name == name {
			return item
		}
	}

	return nil
}
