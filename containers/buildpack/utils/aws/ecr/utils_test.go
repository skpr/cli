package ecr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRegistry(t *testing.T) {
	assert.True(t, IsRegistry("example.ecr.amazon.com"))
}

func TestExtractRegionFromURL(t *testing.T) {
	region, err := extractRegionFromURL("example.ap-southeast-2.aws.amazon.com")
	assert.NoError(t, err)
	assert.Equal(t, "ap-southeast-2", region)

	_, err = extractRegionFromURL("example.nope.aws.amazon.com")
	assert.NotNil(t, err)
}
