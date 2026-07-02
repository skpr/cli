package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getMemoryGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricResourceUsage).GetDate, func(m *pb.MetricResourceUsage) float64 {
		return float64(m.GetMemory())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Memory (MB)", width, computeWindow(points[0].date))
	chart.AddLine("Memory", "\033[34m")

	for _, p := range points {
		chart.PushAt("Memory", p.value, p.date)
	}

	return chart.Render(), nil
}
