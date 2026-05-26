package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// ListEnvironmentsInput is the input schema for the list_environments tool.
type ListEnvironmentsInput struct{}

// EnvironmentSummary is a single row in the list_environments output.
type EnvironmentSummary struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Size       string `json:"size"`
	Routes     string `json:"routes"`
	Phase      string `json:"phase"`
	Production bool   `json:"production"`
}

// ListEnvironmentsOutput is the output schema for the list_environments tool.
type ListEnvironmentsOutput struct {
	Environments []EnvironmentSummary `json:"environments"`
}

// ListEnvironments lists all environments for the current project.
func ListEnvironments(ctx context.Context, req *mcp.CallToolRequest, _ ListEnvironmentsInput) (*mcp.CallToolResult, ListEnvironmentsOutput, error) {
	ctx, skprClient, err := client.New(ctx)
	if err != nil {
		return nil, ListEnvironmentsOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := skprClient.Environment().List(ctx, &pb.EnvironmentListRequest{})
	if err != nil {
		return nil, ListEnvironmentsOutput{}, fmt.Errorf("could not list environments: %w", err)
	}

	out := ListEnvironmentsOutput{
		Environments: make([]EnvironmentSummary, 0, len(resp.Environments)),
	}

	for _, env := range resp.Environments {
		routes := append(env.Ingress.Routes, env.Ingress.Domain)
		out.Environments = append(out.Environments, EnvironmentSummary{
			Name:       env.Name,
			Version:    env.Version,
			Size:       env.Size,
			Routes:     strings.Join(routes, ", "),
			Phase:      env.Phase,
			Production: env.Production,
		})
	}

	return nil, out, nil
}

// GetEnvironmentInput is the input schema for the get_environment tool.
type GetEnvironmentInput struct {
	// Environment is required — no omitempty so the schema marks it required.
	Environment string `json:"environment" jsonschema:"Name of the environment to retrieve"`
}

// GetEnvironmentOutput is the output schema for the get_environment tool.
type GetEnvironmentOutput struct {
	Environment *pb.Environment `json:"environment"`
}

// GetEnvironment retrieves detailed information about a single environment.
func GetEnvironment(ctx context.Context, req *mcp.CallToolRequest, in GetEnvironmentInput) (*mcp.CallToolResult, GetEnvironmentOutput, error) {
	ctx, skprClient, err := client.New(ctx)
	if err != nil {
		return nil, GetEnvironmentOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := skprClient.Environment().Get(ctx, &pb.EnvironmentGetRequest{Name: in.Environment})
	if err != nil {
		return nil, GetEnvironmentOutput{}, fmt.Errorf("could not get environment %q: %w", in.Environment, err)
	}

	return nil, GetEnvironmentOutput{Environment: resp.Environment}, nil
}
