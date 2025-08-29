package list

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
	timeutils "github.com/skpr/cli/internal/time"
)

// Command to list restores.
type Command struct {
	Environment string
	JSON        bool
}

// Row used for formatting list response.
type Row struct {
	RestoreID      string
	Name           string
	BackupID       string
	Phase          string
	StartTime      string
	CompletionTime string
	Duration       string
	Database       string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	resp, err := client.Mysql().RestoreList(ctx, &pb.MysqlRestoreListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return errors.Wrap(err, "restore list failed")
	}

	if len(resp.List) == 0 {
		fmt.Println("No restores found:", "See `skpr restore create` to create a restore")
		return nil
	}

	var rows []Row

	for _, item := range resp.List {
		row := Row{
			Name:      item.Name,
			BackupID:  item.BackupID,
			RestoreID: item.RestoreID,
			Phase:     item.Phase.String(),
			StartTime: item.StartTime,
			Duration:  item.Duration,
			Database:  item.Database,
		}

		if item.StartTime != "" {
			start, err := timeutils.ParseString(item.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time: %w", err)
			}

			row.StartTime = start.Format(time.RFC822Z)
		}

		if item.CompletionTime != "" {
			completion, err := timeutils.ParseString(item.CompletionTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time: %w", err)
			}

			row.CompletionTime = completion.Format(time.RFC822Z)

			row.Duration = item.Duration
		}

		rows = append(rows, row)
	}

	if cmd.JSON {
		data, err := json.Marshal(rows)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	header := []string{
		"Restore ID",
		"Backup ID",
		"Name",
		"Phase",
		"Start Time",
		"Completion Time",
		"Duration",
		"Database",
	}

	var flatRows [][]string

	for _, item := range rows {
		flatRows = append(flatRows, []string{
			item.RestoreID,
			item.Name,
			item.BackupID,
			color.ApplyColorToString(item.Phase),
			item.StartTime,
			item.CompletionTime,
			item.Duration,
			item.Database,
		})
	}

	return table.Print(os.Stdout, header, flatRows)
}
