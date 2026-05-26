package mcp_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	skprmcp "github.com/skpr/cli/internal/command/mcp"
)

// freeAddr finds an available TCP address on localhost by briefly opening a
// listener and then closing it.  The address is returned in "host:port" form.
func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := ln.Addr().String()
	require.NoError(t, ln.Close())
	return addr
}

// waitHealthz polls /healthz until it returns 200 OK or the deadline is
// reached.  It returns an error if the server never becomes healthy.
func waitHealthz(t *testing.T, addr string) error {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	url := "http://" + addr + "/healthz"
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:noctx
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	return context.DeadlineExceeded
}

func TestHTTPTransportLifecycle(t *testing.T) {
	addr := freeAddr(t)

	command := &skprmcp.Command{
		ClientId:   "",
		GitVersion: "v0.0.0-test",
		BuildDate:  "1970-01-01",
		HTTPAddr:   addr,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- command.Run(ctx)
	}()

	// Wait for the HTTP server to be ready.
	require.NoError(t, waitHealthz(t, addr), "server did not become healthy in time")

	// Connect via the Streamable HTTP transport.
	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	session, err := mcpClient.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:             "http://" + addr + "/",
		DisableStandaloneSSE: true,
	}, nil)
	require.NoError(t, err)
	defer session.Close()

	// Collect tool names advertised by the server.
	toolNames := make(map[string]struct{})
	for tool, err := range session.Tools(ctx, nil) {
		require.NoError(t, err)
		toolNames[tool.Name] = struct{}{}
	}

	expected := []string{
		"list_environments",
		"get_environment",
		"mysql_image_list",
		"mysql_image_pull",
		"version",
	}
	for _, name := range expected {
		assert.Contains(t, toolNames, name, "expected tool %q to be registered", name)
	}
	assert.Len(t, toolNames, len(expected))

	// Close the session then cancel the context; the server should shut down cleanly.
	require.NoError(t, session.Close())
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err, "server.Run should return nil on clean shutdown")
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down within 5 seconds")
	}
}

func TestHealthzEndpoint(t *testing.T) {
	addr := freeAddr(t)

	command := &skprmcp.Command{
		ClientId:   "",
		GitVersion: "v0.0.0-test",
		BuildDate:  "1970-01-01",
		HTTPAddr:   addr,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go command.Run(ctx) //nolint:errcheck

	require.NoError(t, waitHealthz(t, addr))

	resp, err := http.Get("http://" + addr + "/healthz") //nolint:noctx
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
