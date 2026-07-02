package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
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

	chart := newChart("Processes", width, computeWindow(points[0].date))
	chart.AddLine("Active", "\033[32m")
	chart.AddLine("Idle", "\033[33m")

	for _, p := range points {
		chart.PushAt("Active", p.active, p.date)
		chart.PushAt("Idle", p.idle, p.date)
	}

	return chart.Render(), nil
}
