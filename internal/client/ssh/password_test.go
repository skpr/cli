package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPassword(t *testing.T) {
	// With session token.
	assert.Equal(t, "foo:bar:baz", getPassword("foo", "bar", "baz"))

	// Without session token.
	assert.Equal(t, "foo:bar", getPassword("foo", "bar", ""))
}
