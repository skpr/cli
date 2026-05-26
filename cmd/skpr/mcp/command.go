package mcp

import (
	"github.com/spf13/cobra"

	"github.com/skpr/cli/containers/docker"
	skprcommand "github.com/skpr/cli/internal/command"
	v1mcp "github.com/skpr/cli/internal/command/mcp"
)

var (
	cmdLong = `Run a Model Context Protocol (MCP) server.

The server exposes a subset of skpr's functionality as MCP tools so that it
can be wired into any MCP-compatible client (Claude Desktop, opencode, etc.).

Available tools
  list_environments  List all environments for the current project.
  get_environment    Get detailed information about a specific environment.
  mysql_image_list   List database images available for an environment.
  mysql_image_pull   Pull a database image to the local Docker daemon.
  version            Return the CLI client and server versions.

Stdio mode (default)
  Reads JSON-RPC from stdin and writes responses to stdout. Use this when
  an MCP client spawns skpr as a child process.

  Example client configuration (opencode / claude_desktop_config.json):

    "skpr": {
      "type": "stdio",
      "command": "skpr",
      "args": ["mcp"]
    }

HTTP mode (--http)
  Starts a Streamable HTTP server on the given address. Use this when running
  skpr as a sidecar container or when multiple clients need to share one server.

  Example — listen on all interfaces, port 8080:
    skpr mcp --http :8080

  Example — localhost only (recommended for local development):
    skpr mcp --http 127.0.0.1:8080

  Example client configuration for HTTP:
    "skpr": {
      "type": "http",
      "url": "http://localhost:8080/"
    }

  Docker sidecar example:
    docker run --rm \
      -v ~/.config/skpr:/root/.config/skpr:ro \
      -p 8080:8080 \
      skpr mcp --http :8080

  Kubernetes sidecar container spec:
    - name: skpr-mcp
      image: ghcr.io/skpr/cli:latest
      args: ["mcp", "--http", ":8080"]
      ports:
        - containerPort: 8080
      volumeMounts:
        - name: skpr-config
          mountPath: /root/.config/skpr
          readOnly: true

  Health probe (liveness / readiness):
    GET /healthz  →  200 OK  body: ok

SECURITY NOTE
  HTTP mode exposes your skpr credentials to anyone who can reach the
  listening port. There is no built-in authentication in this release.
  Bind to localhost or a private network interface, and use your
  infrastructure's network controls to restrict access.`

	// GitVersion is overridden at build time via ldflags (shared with the
	// version command).
	GitVersion string
	// BuildDate is overridden at build time via ldflags.
	BuildDate string
)

// NewCommand creates the cobra.Command for 'skpr mcp'.
func NewCommand(clientId docker.DockerClientId) *cobra.Command {
	command := &v1mcp.Command{
		ClientId:   clientId,
		GitVersion: GitVersion,
		BuildDate:  BuildDate,
	}

	cmd := &cobra.Command{
		Use:                   "mcp",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Run a Model Context Protocol server",
		Long:                  cmdLong,
		GroupID:               skprcommand.GroupOther,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.HTTPAddr, "http", "",
		`Listen address for HTTP transport (e.g. ":8080" or "127.0.0.1:8080"). Empty uses stdio.`)
	cmd.Flags().StringVar(&command.HTTPPath, "path", "/",
		"URL path to mount the MCP handler on (HTTP mode only).")
	cmd.Flags().BoolVar(&command.HTTPStateless, "stateless", false,
		"Disable session tracking (HTTP mode only). Useful behind a stateless load balancer.")
	cmd.Flags().BoolVar(&command.HTTPAllowCrossOriginHosts, "allow-cross-origin-hosts", false,
		"Disable DNS-rebinding protection (HTTP mode only). Only use if you understand the security implications.")

	return cmd
}
