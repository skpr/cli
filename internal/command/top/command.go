package top

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/client"
)

// Command for displaying resource usage metrics.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	var (
		resourceUsage        *pb.ResourceUsageResponse
		responseTimes        *pb.ResponseTimesResponse
		responseCodes        *pb.ResponseCodesResponse
		requests             *pb.RequestsResponse
		cacheRatio           *pb.CacheRatioResponse
		invalidationRequests *pb.InvalidationRequestsResponse
		invalidationPaths    *pb.InvalidationPathsResponse
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		resourceUsage, err = client.Metrics().ResourceUsage(ctx, &pb.ResourceUsageRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		responseTimes, err = client.Metrics().ResponseTimes(ctx, &pb.ResponseTimesRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		responseCodes, err = client.Metrics().ResponseCodes(ctx, &pb.ResponseCodesRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		requests, err = client.Metrics().Requests(ctx, &pb.RequestsRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		cacheRatio, err = client.Metrics().CacheRatio(ctx, &pb.CacheRatioRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		invalidationRequests, err = client.Metrics().InvalidationRequests(ctx, &pb.InvalidationRequestsRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	g.Go(func() error {
		var err error
		invalidationPaths, err = client.Metrics().InvalidationPaths(ctx, &pb.InvalidationPathsRequest{
			Environment: cmd.Environment,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	p := tea.NewProgram(newModel(
		resourceUsage,
		requests,
		cacheRatio,
		invalidationRequests,
		invalidationPaths,
		responseTimes,
		responseCodes,
	), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
