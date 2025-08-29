package create

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command to create a restore.
type Command struct {
	Environment  string
	Backup       string
	DatabaseName string
	Wait         bool
	Force        bool
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
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

	resp, err := client.Mysql().RestoreCreate(ctx, &pb.MysqlRestoreCreateRequest{
		Environment:  cmd.Environment,
		Backup:       cmd.Backup,
		DatabaseName: cmd.DatabaseName,
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

		resp, err := client.Mysql().RestoreGet(ctx, &pb.MysqlRestoreGetRequest{
			ID: resp.ID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get mysql restore")
		}

		switch resp.MysqlRestore.Phase {
		case pb.MysqlRestoreStatus_Completed:
			fmt.Fprintln(os.Stderr, "Restore complete!")
			return nil
		case pb.MysqlRestoreStatus_Failed:
			return fmt.Errorf("the restore failed: the Skpr team has been notified")
		case pb.MysqlRestoreStatus_Unknown:
			return fmt.Errorf("the restore failed for an unknown reason: the Skpr team has been notified")
		}
	}
}
