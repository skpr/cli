package top

import (
	"fmt"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/components/graph/runchart"
)

func getCacheHitRatioGraph(resp *pb.CacheRatioResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricCacheRatio).GetDate, func(m *pb.MetricCacheRatio) float64 {
		return float64(m.GetHit())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Cache Hit Ratio", width, computeWindow(points[0].date), runchart.WithMaxValue(100))
	chart.AddLine("Hit", "\033[32m")

	for _, p := range points {
		chart.PushAt("Hit", p.value, p.date)
	}

	return chart.Render(), nil
}
