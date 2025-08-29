package login

import (
	"context"
	"net/http/httptest"
	"testing"

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
