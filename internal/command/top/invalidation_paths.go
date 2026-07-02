package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getInvalidationPathsGraph(resp *pb.InvalidationPathsResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricInvalidationPaths).GetDate, func(m *pb.MetricInvalidationPaths) float64 {
		return float64(m.GetPaths())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Invalidation Paths", width, computeWindow(points[0].date))
	chart.AddLine("Paths", "\033[35m")

	for _, p := range points {
		chart.PushAt("Paths", p.value, p.date)
	}

	return chart.Render(), nil
}
