package pull

import (
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
)

func TestImageExists(t *testing.T) {
	assert.True(t, imageExists([]*pb.ImageStatus{
		{
			ID: "foo",
		},
	}, "foo"))

	assert.False(t, imageExists([]*pb.ImageStatus{
		{
			ID: "bar",
		},
	}, "foo"))
}
