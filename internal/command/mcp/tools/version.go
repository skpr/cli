package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// VersionInput is the input schema for the version tool.
type VersionInput struct{}

// VersionOutput is the output schema for the version tool.
type VersionOutput struct {
	ClientVersion   string `json:"client_version"`
	ClientBuildDate string `json:"client_build_date"`
	ServerVersion   string `json:"server_version,omitempty"`
	ServerBuildDate string `json:"server_build_date,omitempty"`
}

// NewVersion returns a tool handler for the version tool, embedding build-time
// values supplied by the cobra command layer.
func NewVersion(gitVersion, buildDate string) func(context.Context, *mcp.CallToolRequest, VersionInput) (*mcp.CallToolResult, VersionOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, _ VersionInput) (*mcp.CallToolResult, VersionOutput, error) {
		out := VersionOutput{
			ClientVersion:   gitVersion,
			ClientBuildDate: buildDate,
		}

		// Best-effort: attempt to retrieve the server version.
		ctx, skprClient, err := client.New(ctx)
		if err == nil {
			resp, err := skprClient.Version().Get(ctx, &pb.VersionGetRequest{})
			if err == nil && resp != nil {
				out.ServerVersion = resp.Version
				out.ServerBuildDate = resp.BuildDate
			}
		}

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Client: %s (%s) | Server: %s (%s)",
						or(out.ClientVersion, "unknown"),
						or(out.ClientBuildDate, "unknown"),
						or(out.ServerVersion, "unknown"),
						or(out.ServerBuildDate, "unknown"),
					),
				},
			},
		}

		return result, out, nil
	}
}

func or(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
