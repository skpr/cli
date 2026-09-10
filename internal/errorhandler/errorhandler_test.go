package errorhandler

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/charmbracelet/fang"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/cli/internal/client/credentials"
)

func TestHandler(t *testing.T) {
	t.Run("login required", func(t *testing.T) {
		var buffer bytes.Buffer

		err := fmt.Errorf("%w: oauth2: %q", credentials.ErrLoginRequired, "invalid_grant")

		Handler(&buffer, fang.Styles{}, err)

		assert.Contains(t, buffer.String(), "session has expired")
		assert.Contains(t, buffer.String(), "skpr login")

		// The underlying error is not shown unless a developer asks for it.
		assert.NotContains(t, buffer.String(), "invalid_grant")
	})

	t.Run("login required with debug", func(t *testing.T) {
		t.Setenv(EnvDebug, "true")

		var buffer bytes.Buffer

		err := fmt.Errorf("%w: oauth2: %q", credentials.ErrLoginRequired, "invalid_grant")

		Handler(&buffer, fang.Styles{}, err)

		assert.Contains(t, buffer.String(), "invalid_grant")
	})

	t.Run("all other errors", func(t *testing.T) {
		var buffer bytes.Buffer

		Handler(&buffer, fang.Styles{}, errors.New("environment not found"))

		assert.Contains(t, buffer.String(), "environment not found")
		assert.NotContains(t, buffer.String(), "skpr login")
	})
}
