package project

import (
	"sort"

	"github.com/skpr/api/pb"
)

// ListEnvironmentsByName lists them in printable order.
func ListEnvironmentsByName(project *pb.Project) []string {
	environments := []string{
		project.Environments.Prod,
	}

	nonProd := append([]string(nil), project.Environments.NonProd...)
	sort.Strings(nonProd)
	environments = append(environments, nonProd...)

	return environments
}
