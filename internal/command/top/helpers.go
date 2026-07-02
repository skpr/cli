package top

import (
	"fmt"
	"sort"
	"time"

	"github.com/skpr/cli/internal/components/graph/runchart"
)

type timeValue struct {
	date  time.Time
	value float64
}

func parseSorted[M any](metrics []M, getDate func(M) string, getValue func(M) float64) ([]timeValue, error) {
	points := make([]timeValue, 0, len(metrics))
	for _, m := range metrics {
		date, err := time.Parse(time.RFC3339, getDate(m))
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		points = append(points, timeValue{date: date, value: getValue(m)})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].date.Before(points[j].date)
	})
	return points, nil
}

func computeWindow(earliest time.Time) time.Duration {
	window := time.Since(earliest)
	if window < time.Minute {
		window = time.Minute
	}
	return window
}

func newChart(title string, width int, window time.Duration, extraOpts ...runchart.Option) *runchart.Chart {
	opts := []runchart.Option{
		runchart.WithTitle(title),
		runchart.WithSize(width, 10),
		runchart.WithWindow(window),
		runchart.WithLegend(true),
		runchart.WithMinValue(0),
		runchart.WithTitleColor("\033[38;5;240m"),
	}
	opts = append(opts, extraOpts...)
	return runchart.New(opts...)
}
