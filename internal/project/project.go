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

	sort.Strings(project.Environments.NonProd)
	environments = append(environments, project.Environments.NonProd...)

	return environments
}
