package top

import (
	"fmt"

	"github.com/skpr/api/pb"
)

func getReplicasGraph(resp *pb.ResourceUsageResponse, width int) (string, error) {
	if len(resp.Metrics) == 0 {
		return "", fmt.Errorf("no metrics available")
	}

	points, err := parseSorted(resp.Metrics, (*pb.MetricResourceUsage).GetDate, func(m *pb.MetricResourceUsage) float64 {
		return float64(m.GetReplicas())
	})
	if err != nil {
		return "", err
	}

	chart := newChart("Replicas", width, computeWindow(points[0].date))
	chart.AddLine("Replicas", "\033[36m")

	for _, p := range points {
		chart.PushAt("Replicas", p.value, p.date)
	}

	return chart.Render(), nil
}
