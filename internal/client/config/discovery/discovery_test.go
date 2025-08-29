package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupProjectDir(t *testing.T) {
	projectRootCmd = func() (string, error) {
		return "foo", nil
	}

	disco, err := New()
	assert.NoError(t, err)

	config, err := disco.Config()
	assert.NoError(t, err)

	assert.Equal(t, "foo/.skpr/config.yml", config)

}
