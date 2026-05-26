package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/skpr/cli/containers/docker"
)

// Command runs the MCP server over stdio or HTTP.
type Command struct {
	ClientId   docker.DockerClientId
	GitVersion string
	BuildDate  string

	// HTTP transport options. HTTPAddr being non-empty selects HTTP mode.
	// In HTTP mode stdout is free for logging; stdio mode keeps stdout
	// exclusively for JSON-RPC frames.
	HTTPAddr                  string // e.g. ":8080" or "127.0.0.1:8080"
	HTTPPath                  string // URL path to mount the MCP handler on (default "/")
	HTTPStateless             bool   // maps to StreamableHTTPOptions.Stateless
	HTTPAllowCrossOriginHosts bool   // maps to StreamableHTTPOptions.DisableLocalhostProtection
}

// Run starts the MCP server and blocks until the client disconnects or the
// context is cancelled.
//
// When HTTPAddr is empty the server runs over stdio (stdin/stdout). All
// JSON-RPC framing is on stdout; incidental output goes to stderr.
//
// When HTTPAddr is non-empty the server listens for Streamable HTTP
// connections on the given address. stdout is available for logging.
func (cmd *Command) Run(ctx context.Context) error {
	srv := Build(cmd.ClientId, cmd.GitVersion, cmd.BuildDate)
	if cmd.HTTPAddr == "" {
		return srv.Run(ctx, &mcp.StdioTransport{})
	}
	return cmd.runHTTP(ctx, srv)
}

// runHTTP starts an HTTP server that wraps the MCP server using the
// Streamable HTTP transport defined by the MCP spec.
func (cmd *Command) runHTTP(ctx context.Context, srv *mcp.Server) error {
	path := cmd.HTTPPath
	if path == "" {
		path = "/"
	}

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return srv },
		&mcp.StreamableHTTPOptions{
			Stateless:                  cmd.HTTPStateless,
			DisableLocalhostProtection: cmd.HTTPAllowCrossOriginHosts,
			Logger:                     slog.New(slog.NewTextHandler(os.Stderr, nil)),
		},
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	httpSrv := &http.Server{
		Addr:              cmd.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("MCP server listening", "addr", cmd.HTTPAddr, "path", path)
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
