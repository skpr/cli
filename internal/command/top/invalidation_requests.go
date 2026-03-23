package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/components/graph/runchart"
)

func getInvalidationRequestsGraph(resp *pb.InvalidationRequestsResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	type dataPoint struct {
		date     time.Time
		requests float64
	}

	points := make([]dataPoint, 0, len(resp.Metrics))

	for _, metric := range resp.Metrics {
		date, err := time.Parse(time.RFC3339, metric.GetDate())
		if err != nil {
			return "", fmt.Errorf("failed to parse date: %w", err)
		}

		points = append(points, dataPoint{
			date:     date,
			requests: float64(metric.GetRequests()),
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].date.Before(points[j].date)
	})

	earliest := points[0].date

	window := time.Since(earliest)
	if window < time.Minute {
		window = time.Minute
	}

	chart := runchart.New(
		runchart.WithTitle("Invalidation Requests"),
		runchart.WithSize(width, 10),
		runchart.WithWindow(window),
		runchart.WithLegend(true),
		runchart.WithMinValue(0),
		runchart.WithTitleColor("\033[38;5;240m"),
	)
	chart.AddLine("Requests", "\033[33m") // yellow

	for _, p := range points {
		chart.PushAt("Requests", p.requests, p.date)
	}

	return chart.Render(), nil
}
