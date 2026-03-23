package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/components/graph/runchart"
)

func getResponseCodesGraph(resp *pb.ResponseCodesResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	type dataPoint struct {
		date        time.Time
		successful  float64
		client      float64
		server      float64
		redirection float64
	}

	points := make([]dataPoint, 0, len(resp.Metrics))

	for _, metric := range resp.Metrics {
		date, err := time.Parse(time.RFC3339, metric.GetDate())
		if err != nil {
			return "", fmt.Errorf("failed to parse date: %w", err)
		}

		points = append(points, dataPoint{
			date:        date,
			successful:  float64(metric.GetSuccessful()),
			client:      float64(metric.GetClient()),
			server:      float64(metric.GetServer()),
			redirection: float64(metric.GetRedirection()),
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
		runchart.WithTitle("Response Codes"),
		runchart.WithSize(width, 10),
		runchart.WithWindow(window),
		runchart.WithLegend(true),
		runchart.WithMinValue(0),
		runchart.WithTitleColor("\033[38;5;240m"),
	)
	chart.AddLine("2xx", "\033[32m") // green
	chart.AddLine("4xx", "\033[33m") // yellow
	chart.AddLine("5xx", "\033[31m") // red
	chart.AddLine("3xx", "\033[36m") // cyan

	for _, p := range points {
		chart.PushAt("2xx", p.successful, p.date)
		chart.PushAt("4xx", p.client, p.date)
		chart.PushAt("5xx", p.server, p.date)
		chart.PushAt("3xx", p.redirection, p.date)
	}

	return chart.Render(), nil
}
