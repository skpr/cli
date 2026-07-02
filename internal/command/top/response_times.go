package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
)

func getResponseTimesGraph(resp *pb.ResponseTimesResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	type dataPoint struct {
		date    time.Time
		average float64
		p95     float64
		p99     float64
	}

	points := make([]dataPoint, 0, len(resp.Metrics))
	for _, metric := range resp.Metrics {
		date, err := time.Parse(time.RFC3339, metric.GetDate())
		if err != nil {
			return "", fmt.Errorf("failed to parse date: %w", err)
		}
		points = append(points, dataPoint{
			date:    date,
			average: float64(metric.GetAverage()),
			p95:     float64(metric.GetP95()),
			p99:     float64(metric.GetP99()),
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].date.Before(points[j].date)
	})

	chart := newChart("Response Times (ms)", width, computeWindow(points[0].date))
	chart.AddLine("Average", "\033[32m")
	chart.AddLine("P95", "\033[33m")
	chart.AddLine("P99", "\033[31m")

	for _, p := range points {
		chart.PushAt("Average", p.average, p.date)
		chart.PushAt("P95", p.p95, p.date)
		chart.PushAt("P99", p.p99, p.date)
	}

	return chart.Render(), nil
}
