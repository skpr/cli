package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getCPUGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricResourceUsage).GetDate, func(m *pb.MetricResourceUsage) float64 {
		return float64(m.GetCPU())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("CPU %", width, computeWindow(points[0].date))
	chart.AddLine("CPU", "\033[32m")

	for _, p := range points {
		chart.PushAt("CPU", p.value, p.date)
	}

	return chart.Render(), nil
}
