package tools_test

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	skprmcp "github.com/skpr/cli/internal/command/mcp"
)

// buildTestServer builds the MCP server with a dummy docker client ID and
// empty build info so it can be used in tests without real credentials.
func buildTestServer() *mcp.Server {
	return skprmcp.Build("", "v0.0.0-test", "1970-01-01")
}

func TestToolRegistration(t *testing.T) {
	srv := buildTestServer()

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	_, err := srv.Connect(ctx, st, nil)
	require.NoError(t, err)

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	session, err := client.Connect(ctx, ct, nil)
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

	assert.Len(t, toolNames, len(expected), "unexpected number of registered tools")
}

func TestMysqlImagePullInputSchema(t *testing.T) {
	srv := buildTestServer()

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	_, err := srv.Connect(ctx, st, nil)
	require.NoError(t, err)

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	session, err := mcpClient.Connect(ctx, ct, nil)
	require.NoError(t, err)
	defer session.Close()

	// Find the mysql_image_pull tool and verify its input schema marks
	// "environment" as required.
	for tool, err := range session.Tools(ctx, nil) {
		require.NoError(t, err)
		if tool.Name != "mysql_image_pull" {
			continue
		}

		require.NotNil(t, tool.InputSchema, "mysql_image_pull must have an input schema")

		// InputSchema arrives as map[string]any over the wire from the client.
		schema, ok := tool.InputSchema.(map[string]any)
		require.True(t, ok, "InputSchema should be a map[string]any")

		rawRequired, ok := schema["required"]
		require.True(t, ok, "mysql_image_pull input schema must have a 'required' key")

		required, ok := rawRequired.([]any)
		require.True(t, ok, "'required' must be a []any")

		requiredStrings := make([]string, len(required))
		for i, r := range required {
			requiredStrings[i], _ = r.(string)
		}

		assert.Contains(t, requiredStrings, "environment", "environment must be required in mysql_image_pull input schema")
		return
	}

	t.Fatal("mysql_image_pull tool not found")
}
