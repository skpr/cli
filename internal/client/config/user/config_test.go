package user

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteConfig(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".skpr", "config.yml")

	configFile := ConfigFile{Path: filePath}

	config := Config{
		Aliases: Aliases{"foo": "bar baz", "whiz": "whang woo"},
	}

	err := configFile.Write(config)
	assert.NoError(t, err)

	assert.FileExists(t, filePath)

	newCfg, err := configFile.Read()
	assert.NoError(t, err)

	aliases := newCfg.Aliases

	assert.Equal(t, "bar baz", aliases["foo"])
	assert.Equal(t, "whang woo", aliases["whiz"])
}

func TestConfigExists(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".skpr", "config.yml")

	configFile := ConfigFile{Path: filePath}
	exists, err := configFile.Exists()
	assert.NoError(t, err)
	assert.False(t, exists)
}
