package buildpack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReference(t *testing.T) {

	image := Image{
		Name: "foo/bar",
		Tag:  "baz",
	}

	assert.Equal(t, "foo/bar:baz", image.Reference())
}
