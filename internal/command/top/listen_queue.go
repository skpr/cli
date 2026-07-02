package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getListenQueueGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricResourceUsage).GetDate, func(m *pb.MetricResourceUsage) float64 {
		return float64(m.GetListenQueue())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Listen Queue", width, computeWindow(points[0].date))
	chart.AddLine("Queue", "\033[31m")

	for _, p := range points {
		chart.PushAt("Queue", p.value, p.date)
	}

	return chart.Render(), nil
}
