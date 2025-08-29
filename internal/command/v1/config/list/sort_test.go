package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/api/pb"
)

func TestSortConfigList(t *testing.T) {
	list := []*pb.Config{
		{
			Key:   "second",
			Value: "yyyyyyyyyyyyyy",
		},
		{
			Key:   "first",
			Value: "xxxxxxxxxxxxxx",
		},
	}

	want := []*pb.Config{
		{
			Key:   "first",
			Value: "xxxxxxxxxxxxxx",
		},
		{
			Key:   "second",
			Value: "yyyyyyyyyyyyyy",
		},
	}

	assert.Equal(t, want, sortConfig(list))
}
