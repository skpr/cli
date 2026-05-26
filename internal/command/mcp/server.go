package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/skpr/cli/containers/docker"
	"github.com/skpr/cli/internal/command/mcp/tools"
)

// Build creates a new MCP server with all registered tools.
func Build(clientId docker.DockerClientId, gitVersion, buildDate string) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "skpr",
		Version: gitVersion,
	}, nil)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_environments",
		Description: "List all environments for the current project and their status.",
	}, tools.ListEnvironments)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_environment",
		Description: "Get detailed information about a specific environment.",
	}, tools.GetEnvironment)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "mysql_image_list",
		Description: "List database images available for an environment.",
	}, tools.MysqlImageList)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "mysql_image_pull",
		Description: "Pull a database image for an environment to the local Docker daemon.",
	}, tools.NewMysqlImagePull(clientId))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "version",
		Description: "Return the CLI client version and (if reachable) the server version.",
	}, tools.NewVersion(gitVersion, buildDate))

	return srv
}
