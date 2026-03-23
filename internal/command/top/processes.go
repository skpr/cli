package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/components/graph/runchart"
)

func getProcessesGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	type dataPoint struct {
		date   time.Time
		active float64
		idle   float64
	}

	points := make([]dataPoint, 0, len(resp.Metrics))

	for _, metric := range resp.Metrics {
		date, err := time.Parse(time.RFC3339, metric.GetDate())
		if err != nil {
			return "", fmt.Errorf("failed to parse date: %w", err)
		}

		points = append(points, dataPoint{
			date:   date,
			active: float64(metric.GetActiveProcesses()),
			idle:   float64(metric.GetIdleProcesses()),
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
		runchart.WithTitle("Processes"),
		runchart.WithSize(width, 10),
		runchart.WithWindow(window),
		runchart.WithLegend(true),
		runchart.WithMinValue(0),
		runchart.WithTitleColor("\033[38;5;240m"),
	)
	chart.AddLine("Active", "\033[32m") // green
	chart.AddLine("Idle", "\033[33m")   // yellow

	for _, p := range points {
		chart.PushAt("Active", p.active, p.date)
		chart.PushAt("Idle", p.idle, p.date)
	}

	return chart.Render(), nil
}
