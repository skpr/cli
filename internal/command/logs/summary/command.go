package summary

import (
	"context"
	"fmt"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/color"
	logshared "github.com/skpr/cli/internal/command/logs"
)

// Command runs an AI summarisation over a log window.
type Command struct {
	Environment string
	Streams     []string
	Since       time.Duration
	From        string
	To          string
	Contains    []string
	Exclude     []string
	Prompt      string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	if cmd.Prompt == "" {
		return fmt.Errorf("a prompt must be provided")
	}

	filter, err := logshared.BuildFilter(logshared.FilterParams{
		Environment: cmd.Environment,
		Streams:     cmd.Streams,
		Since:       cmd.Since,
		From:        cmd.From,
		To:          cmd.To,
		Contains:    cmd.Contains,
		Exclude:     cmd.Exclude,
	})
	if err != nil {
		return err
	}

	ctx, c, err := client.New(ctx)
	if err != nil {
		return err
	}

	resp, err := c.Logs().Summarise(ctx, &pb.LogSummariseRequest{
		Filter: filter,
		Prompt: cmd.Prompt,
	})
	if err != nil {
		return fmt.Errorf("failed to summarise logs: %w", err)
	}

	fmt.Println(orangeBold("▌ Overview"))
	fmt.Printf("  %s\n", resp.GetOverview())

	if bullets := resp.GetBullets(); len(bullets) > 0 {
		fmt.Println()
		fmt.Println(orangeBold("▌ Notable"))
		for _, b := range bullets {
			fmt.Printf("  • %s\n", b)
		}
	}

	if actions := resp.GetSuggestedActions(); len(actions) > 0 {
		fmt.Println()
		fmt.Println(orangeBold("▌ Suggested actions"))
		for _, a := range actions {
			fmt.Printf("  • %s\n", a)
		}
	}

	fmt.Println()

	return nil
}

func orangeBold(s string) string {
	return gchalk.WithHex(color.HexOrange).Bold(s)
}
