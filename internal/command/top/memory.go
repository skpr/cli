package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/components/graph/runchart"
)

func getMemoryGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	// Parse all dates to determine the time window.
	type dataPoint struct {
		date   time.Time
		memory float64
	}

	points := make([]dataPoint, 0, len(resp.Metrics))

	for _, metric := range resp.Metrics {
		date, err := time.Parse(time.RFC3339, metric.GetDate())
		if err != nil {
			return "", fmt.Errorf("failed to parse date: %w", err)
		}

		points = append(points, dataPoint{
			date:   date,
			memory: float64(metric.GetMemory()),
		})
	}

	// Sort by date so the chart draws lines in chronological order.
	sort.Slice(points, func(i, j int) bool {
		return points[i].date.Before(points[j].date)
	})

	// Compute the window from the earliest data point to now.
	earliest := points[0].date

	window := time.Since(earliest)
	if window < time.Minute {
		window = time.Minute
	}

	memoryChart := runchart.New(
		runchart.WithTitle("Memory (MB)"),
		runchart.WithSize(width, 10),
		runchart.WithWindow(window),
		runchart.WithLegend(true),
		runchart.WithMinValue(0),
		runchart.WithTitleColor("\033[38;5;240m"),
	)
	memoryChart.AddLine("Memory", "\033[34m") // blue

	for _, p := range points {
		memoryChart.PushAt("Memory", p.memory, p.date)
	}

	return memoryChart.Render(), nil
}
