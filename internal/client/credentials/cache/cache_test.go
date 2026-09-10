package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	// Redirect the cache directory. We set both because os.UserCacheDir uses
	// XDG_CACHE_HOME on Linux and HOME on macOS.
	directory := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", directory)
	t.Setenv("HOME", directory)

	_, found, err := Get("example.com")
	require.NoError(t, err)
	assert.False(t, found, "credentials should not exist before they have been set")

	err = Set("example.com", Credentials{Token: Token{Refresh: "test-refresh-token"}})
	require.NoError(t, err)

	credentials, found, err := Get("example.com")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "test-refresh-token", credentials.Token.Refresh)

	// These credentials contain a refresh token, so they should only be
	// readable by the current user.
	path, err := getFile("example.com")
	require.NoError(t, err)

	file, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(FilePerms), file.Mode().Perm())

	cacheDirectory, err := getDirectory()
	require.NoError(t, err)

	dir, err := os.Stat(cacheDirectory)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(DirectoryPerms), dir.Mode().Perm())

	err = Delete("example.com")
	require.NoError(t, err)

	_, found, err = Get("example.com")
	require.NoError(t, err)
	assert.False(t, found)

	// Deleting credentials which are already gone is not an error.
	assert.NoError(t, Delete("example.com"))
}
