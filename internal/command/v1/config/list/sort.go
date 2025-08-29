package list

import (
	"sort"

	"github.com/skpr/api/pb"
)

// Helper function to sort a list of configs.
func sortConfig(list []*pb.Config) []*pb.Config {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Key < list[j].Key
	})

	return list
}
