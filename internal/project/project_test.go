package project

import (
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	item := pb.Project{
		Environments: &pb.ProjectEnvironments{
			Prod: "prod",
			NonProd: []string{
				"dev",
				"training",
				"stg",
				"bvt",
			},
		},
	}

	out := ListEnvironmentsByName(&item)
	expected := []string{"prod", "bvt", "dev", "stg", "training"}

	assert.Equal(t, expected, out)
}
