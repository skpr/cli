package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getInvalidationRequestsGraph(resp *pb.InvalidationRequestsResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricInvalidationRequests).GetDate, func(m *pb.MetricInvalidationRequests) float64 {
		return float64(m.GetRequests())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Invalidation Requests", width, computeWindow(points[0].date))
	chart.AddLine("Requests", "\033[33m")

	for _, p := range points {
		chart.PushAt("Requests", p.value, p.date)
	}

	return chart.Render(), nil
}
