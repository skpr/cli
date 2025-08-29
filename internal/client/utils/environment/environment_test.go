package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/api/pb"
)

func TestContains(t *testing.T) {
	environments := []*pb.Environment{
		{
			Name: "dev",
		},
		{
			Name: "stg",
		},
	}

	assert.True(t, Contains("dev", environments))
	assert.False(t, Contains("staging", environments))
}
