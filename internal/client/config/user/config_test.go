package command

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteConfig(t *testing.T) {
	filename := filepath.Join(t.TempDir(), ".skpr", "config.yml")
	configFile := NewConfigFile(filename)
	config := Config{
		Aliases: Aliases{"foo": "bar baz", "whiz": "whang woo"},
	}
	err := configFile.Write(config)
	assert.NoError(t, err)

	assert.FileExists(t, filename)

	newCfg, err := configFile.Read()
	assert.NoError(t, err)

	aliases := newCfg.Aliases

	assert.Equal(t, "bar baz", aliases["foo"])
	assert.Equal(t, "whang woo", aliases["whiz"])
}

func TestConfigExists(t *testing.T) {
	filename := filepath.Join(t.TempDir(), ".skpr", "config.yml")
	configFile := NewConfigFile(filename)
	exists, err := configFile.Exists()
	assert.NoError(t, err)
	assert.False(t, exists)
}
