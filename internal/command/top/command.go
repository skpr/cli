package top

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/client"
)

// Command for displaying resource usage metrics.
type Command struct {
	Environment string
	Refresh     time.Duration
}

func fetchAllMetrics(ctx context.Context, c *client.Client, env string) (metricsMsg, error) {
	var msg metricsMsg

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		msg.resourceUsage, err = c.Metrics().ResourceUsage(ctx, &pb.ResourceUsageRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.responseTimes, err = c.Metrics().ResponseTimes(ctx, &pb.ResponseTimesRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.responseCodes, err = c.Metrics().ResponseCodes(ctx, &pb.ResponseCodesRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.requests, err = c.Metrics().Requests(ctx, &pb.RequestsRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.cacheRatio, err = c.Metrics().CacheRatio(ctx, &pb.CacheRatioRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.invalidationRequests, err = c.Metrics().InvalidationRequests(ctx, &pb.InvalidationRequestsRequest{
			Environment: env,
		})
		return err
	})

	g.Go(func() error {
		var err error
		msg.invalidationPaths, err = c.Metrics().InvalidationPaths(ctx, &pb.InvalidationPathsRequest{
			Environment: env,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return metricsMsg{}, err
	}

	return msg, nil
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, c, err := client.New(ctx)
	if err != nil {
		return err
	}

	initial, err := fetchAllMetrics(ctx, c, cmd.Environment)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	fetchMetrics := func() tea.Msg {
		msg, err := fetchAllMetrics(ctx, c, cmd.Environment)
		if err != nil {
			return metricsErrMsg{err: err}
		}
		return msg
	}

	p := tea.NewProgram(newModel(initial, cmd.Refresh, fetchMetrics), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
