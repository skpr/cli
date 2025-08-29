package aliases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/cli/internal/client/config/command"
)

func TestExpandAliases(t *testing.T) {
	args := []string{"mp", "dev"}
	aliases := command.Aliases{
		"mp": "mysql image pull",
	}
	found, newArgs, err := Expand(args, aliases)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, []string{"mysql", "image", "pull", "dev"}, newArgs)
}

func TestExpandAliasesWithPlaceholder(t *testing.T) {
	aliases := command.Aliases{
		"mp": "mysql image pull $1",
		"fs": "rsync $1:/data/app/sites/default/files $2",
	}
	args := []string{"mp", "dev"}
	found, newArgs, err := Expand(args, aliases)
	assert.NoError(t, err)
	assert.True(t, found)

	assert.Equal(t, []string{"mysql", "image", "pull", "dev"}, newArgs)

	args = []string{"fs", "dev", "."}
	found, newArgs, err = Expand(args, aliases)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, []string{"rsync", "dev:/data/app/sites/default/files", "."}, newArgs)
}

func TestExpandNoArgs(t *testing.T) {
	args := []string{}
	aliases := command.Aliases{}
	found, newArgs, err := Expand(args, aliases)
	assert.NoError(t, err)
	assert.Len(t, newArgs, 0)
	assert.False(t, found)
}
