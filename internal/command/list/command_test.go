package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/api/pb"
)

func TestSort(t *testing.T) {
	list := []*pb.Environment{
		{
			Name:       "dev",
			Production: false,
		},
		{
			Name:       "prod",
			Production: true,
		},
		{
			Name:       "stg",
			Production: false,
		},
	}

	sortEnvs(list)

	assert.Equal(t, "dev", list[0].Name)
	assert.Equal(t, "stg", list[1].Name)
	assert.Equal(t, "prod", list[2].Name)
}
