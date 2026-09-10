package login

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleCallback(t *testing.T) {
	server := NewServer("http://localhost:8080")

	_, shutdown := context.WithCancel(context.Background())

	handler := server.handleLoginCallback(shutdown)

	request := httptest.NewRequest("GET", "/?code=123456&state=789&error=ERROR_TEST&error_description=ERROR_DESCRIPTION", nil)
	response := httptest.NewRecorder()

	// Execute the callback handler.
	handler(response, request)

	// Check that the response fields were set.
	assert.Equal(t, "123456", server.Response.Code)
	assert.Equal(t, "789", server.Response.State)
	assert.Equal(t, "ERROR_TEST", server.Response.Error)
	assert.Equal(t, "ERROR_DESCRIPTION", server.Response.ErrorDescription)
}

func TestRunWithInvalidCallback(t *testing.T) {
	tests := []struct {
		name     string
		callback string
	}{
		{
			name:     "callback not provided",
			callback: "",
		},
		{
			name:     "callback without a port",
			callback: "http://localhost",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := NewServer(test.callback)

			_, ready := context.WithCancel(context.Background())
			defer ready()

			// This must fail immediately. Previously the command would wait for
			// a callback which was never going to arrive.
			_, err := server.Run(context.Background(), ready)
			assert.ErrorContains(t, err, "callback URL must include a host and a port")
		})
	}
}

func TestRunWithoutCallback(t *testing.T) {
	server := NewServer("http://localhost:11219")

	// The identity provider does not redirect back to us unless our callback
	// URL has been registered, so we must not wait forever.
	server.CallbackTimeout = 250 * time.Millisecond

	_, ready := context.WithCancel(context.Background())
	defer ready()

	_, err := server.Run(context.Background(), ready)
	assert.ErrorIs(t, err, ErrCallbackTimeout)
}
