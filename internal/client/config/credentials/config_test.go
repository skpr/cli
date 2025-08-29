package credentials

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	file := filepath.Join(os.TempDir(), "credentials.yml")
	config := NewConfig(file)
	clusterKey := "test.cluster.skpr.io"
	err := config.AddOrUpdate(clusterKey, Cluster{
		AWS: AWS{
			Profile: "test.profile",
		},
	})
	assert.NoError(t, err)

	assert.FileExists(t, file)

	cluster, err := config.Get(clusterKey)
	assert.NoError(t, err)
	assert.Equal(t, "test.profile", cluster.AWS.Profile)

}
