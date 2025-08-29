package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/table"
)

// Command to list all restores.
type Command struct {
	Environment string
	JSON        bool
}

// Row used for formatting list response.
type Row struct {
	Name      string   `header:"name"`
	Backup    string   `header:"backup"`
	Phase     string   `header:"phase"`
	StartTime string   `header:"start time"`
	Duration  string   `header:"duration"`
	Databases []string `header:"databases"`
	Volumes   []string `header:"volumes"`
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	resp, err := client.Restore().List(ctx, &pb.RestoreListRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	if len(resp.List) == 0 {
		fmt.Println("No restores found:", "See `skpr restore create` to create a restore")
		return nil
	}

	var rows []Row

	for _, item := range resp.List {
		databaseList := []string{
			"n/a",
		}

		if len(item.Databases) > 0 {
			databaseList = item.Databases
		}

		volumeList := []string{
			"n/a",
		}

		if len(item.Volumes) > 0 {
			volumeList = item.Volumes
		}

		row := Row{
			Name:      item.Name,
			Backup:    item.Backup,
			Phase:     item.Phase.String(),
			StartTime: item.StartTime,
			Duration:  item.Duration,
			Databases: databaseList,
			Volumes:   volumeList,
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
		"Name",
		"Backup ID",
		"Phase",
		"Start Time",
		"Duration",
		"Databases",
		"Volumes",
	}

	var flatRows [][]string

	for _, item := range rows {
		flatRows = append(flatRows, []string{
			item.Name,
			item.Backup,
			color.ApplyColorToString(item.Phase),
			item.StartTime,
			item.Duration,
			strings.Join(item.Databases, ","),
			strings.Join(item.Volumes, ","),
		})
	}

	return table.Print(os.Stdout, header, flatRows)
}
