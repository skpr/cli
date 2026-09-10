package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	assert.NoError(t, Wait(context.Background(), server.URL, time.Second))

	// A server which never replies returns after the timeout instead of
	// waiting forever.
	err := Wait(context.Background(), "http://127.0.0.1:1/readyz", 100*time.Millisecond)
	assert.ErrorContains(t, err, "server did not reply after")

	// Waiting stops when the caller gives up.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.ErrorIs(t, Wait(ctx, "http://127.0.0.1:1/readyz", time.Minute), context.Canceled)
}
