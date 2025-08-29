package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindDockerfiles(t *testing.T) {
	dockerfiles, err := FindDockerfiles("testdata/.skpr")
	assert.NoError(t, err)

	assert.Len(t, dockerfiles, 4)
	assert.Equal(t, "testdata/.skpr/compile/Dockerfile", dockerfiles["compile"])
	assert.Equal(t, "testdata/.skpr/cli/Dockerfile", dockerfiles["cli"])
	assert.Equal(t, "testdata/.skpr/app/Dockerfile", dockerfiles["app"])
	assert.Equal(t, "testdata/.skpr/web/Dockerfile", dockerfiles["web"])
}

func TestFindLegacyDockerfiles(t *testing.T) {
	dockerfiles, err := FindDockerfiles("testdata/legacy")
	assert.NoError(t, err)

	assert.Len(t, dockerfiles, 4)
	assert.Equal(t, "testdata/legacy/compile.dockerfile", dockerfiles["compile"])
	assert.Equal(t, "testdata/legacy/cli.dockerfile", dockerfiles["cli"])
	assert.Equal(t, "testdata/legacy/nginx.dockerfile", dockerfiles["nginx"])
	assert.Equal(t, "testdata/legacy/php.dockerfile", dockerfiles["php"])
}
