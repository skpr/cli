package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	"github.com/skpr/cli/internal/components/tooltip"
	"github.com/skpr/cli/internal/table"
)

const (
	// MaxValueLength to be applied when listing values.
	MaxValueLength = 40
)

// Command for listing config.
type Command struct {
	Environment string
	JSON        bool
	FilterType  string
	ShowSecrets bool
	Wide        bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Config().List(ctx, &pb.ConfigListRequest{
		Name:        cmd.Environment,
		FilterType:  getFilterType(cmd.FilterType),
		ShowSecrets: cmd.ShowSecrets,
	})
	if err != nil {
		return err
	}

	if cmd.JSON {
		items := map[string]string{}
		for _, item := range resp.List {
			items[item.Key] = item.Value
		}
		data, err := json.Marshal(items)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return Print(os.Stdout, resp.List, cmd.Environment, MaxValueLength, cmd.Wide)
}

// Helper function to sort a list of configs.
func sortConfigList(list []*pb.Config) []*pb.Config {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Key < list[j].Key
	})

	return list
}

// Row which can be....
type Row struct {
	Key   string `header:"key"`
	Value string `header:"value"`
	Type  string `header:"type"`
}

// Print the table...
func Print(w io.Writer, list []*pb.Config, environment string, trimLen int, wide bool) error {
	header := []string{
		"Key",
		"Value",
		"Type",
	}

	var rows [][]string

	// A list of all the trimmed config keys.
	var trimmed bool

	for _, item := range sortConfigList(list) {
		if len(item.Value) > trimLen && !wide {
			item.Value = fmt.Sprintf("%s...", item.Value[:trimLen])
			trimmed = true
		}

		rows = append(rows, []string{
			item.Key,
			item.Value,
			color.ApplyColorToString(item.Type.String()),
		})
	}

	err := table.Print(w, header, rows)
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	var tooltipText = "Show secret values using this command:\n\n"
	tooltipText = fmt.Sprintf("%s$ skpr config list %s --show-secrets\n", tooltipText, environment)

	if trimmed {
		tooltipText = fmt.Sprintf("%s\nValues have been trimmed. See the full value using this command:\n\n", tooltipText)
		tooltipText = fmt.Sprintf("%s$ skpr config list %s --wide\n", tooltipText, environment)
	}

	err = tooltip.Render(w, tooltipText)
	if err != nil {
		return fmt.Errorf("failed to render tooltip: %w", err)
	}

	return nil
}

// Helper function to convert a string to a config type.
func getFilterType(name string) pb.ConfigType {
	name = strings.ToLower(name)
	if name == strings.ToLower(pb.ConfigType_System.String()) {
		return pb.ConfigType_System
	}

	if name == strings.ToLower(pb.ConfigType_User.String()) {
		return pb.ConfigType_User
	}

	if name == strings.ToLower(pb.ConfigType_Overridden.String()) {
		return pb.ConfigType_Overridden
	}

	return pb.ConfigType_None
}
