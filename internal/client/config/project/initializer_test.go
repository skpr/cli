package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skpr/cli/internal/random"
)

func TestInitialize(t *testing.T) {

	assert := require.New(t)

	dir, err := os.MkdirTemp(os.TempDir(), random.String(8))
	assert.NoError(err)

	sub, err := fs.Sub(resources, "resources")

	f, err := sub.Open("config.yml")
	assert.NoError(err)

	assert.NotEmpty(f)

	initializer := NewInitializer(dir, "test.cluster", "test-project")
	err = initializer.Initialize()
	assert.NoError(err)

	assert.FileExists(filepath.Join(dir, "config.yml"))
	dat, err := os.ReadFile(filepath.Join(dir, "config.yml"))
	assert.NoError(err)
	assert.Contains(string(dat), "test.cluster")
	assert.Contains(string(dat), "test-project")
}
