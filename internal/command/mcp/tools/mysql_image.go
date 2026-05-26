package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/containers/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/containers/docker"
	"github.com/skpr/cli/containers/docker/types"
	"github.com/skpr/cli/internal/client"
	timeutils "github.com/skpr/cli/internal/time"
)

// MysqlImageListInput is the input schema for the mysql_image_list tool.
type MysqlImageListInput struct {
	// Environment is required — no omitempty so the schema marks it required.
	Environment string `json:"environment" jsonschema:"Name of the environment"`
}

// MysqlImageRow is a single row in the mysql_image_list output.
type MysqlImageRow struct {
	ID             string   `json:"id"`
	Phase          string   `json:"phase"`
	StartTime      string   `json:"start_time,omitempty"`
	CompletionTime string   `json:"completion_time,omitempty"`
	Duration       string   `json:"duration,omitempty"`
	Tags           []string `json:"tags"`
}

// MysqlImageListOutput is the output schema for the mysql_image_list tool.
type MysqlImageListOutput struct {
	Images []MysqlImageRow `json:"images"`
}

// MysqlImageList lists database images for the given environment.
func MysqlImageList(ctx context.Context, req *mcp.CallToolRequest, in MysqlImageListInput) (*mcp.CallToolResult, MysqlImageListOutput, error) {
	ctx, skprClient, err := client.New(ctx)
	if err != nil {
		return nil, MysqlImageListOutput{}, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := skprClient.Mysql().ImageList(ctx, &pb.ImageListRequest{
		Environment: in.Environment,
	})
	if err != nil {
		return nil, MysqlImageListOutput{}, fmt.Errorf("image list failed: %w", err)
	}

	out := MysqlImageListOutput{
		Images: make([]MysqlImageRow, 0, len(resp.List)),
	}

	for _, item := range resp.List {
		row := MysqlImageRow{
			ID:    item.ID,
			Phase: item.Phase.String(),
			Tags:  item.Tags,
		}

		if item.StartTime != "" {
			start, err := timeutils.ParseString(item.StartTime)
			if err != nil {
				return nil, MysqlImageListOutput{}, fmt.Errorf("failed to parse start time: %w", err)
			}
			row.StartTime = start.Format(time.RFC822Z)
		}

		if item.CompletionTime != "" {
			completion, err := timeutils.ParseString(item.CompletionTime)
			if err != nil {
				return nil, MysqlImageListOutput{}, fmt.Errorf("failed to parse completion time: %w", err)
			}
			row.CompletionTime = completion.Format(time.RFC822Z)
			row.Duration = item.Duration
		}

		out.Images = append(out.Images, row)
	}

	return nil, out, nil
}

// MysqlImagePullInput is the input schema for the mysql_image_pull tool.
type MysqlImagePullInput struct {
	// Environment is required — no omitempty so the schema marks it required.
	Environment string   `json:"environment"         jsonschema:"Name of the environment"`
	Databases   []string `json:"databases,omitempty" jsonschema:"Database names to pull (defaults to [default])"`
	ID          string   `json:"id,omitempty"        jsonschema:"Specific image ID to pull (overrides Databases)"`
}

// MysqlImagePullEntry is a single pulled image result.
type MysqlImagePullEntry struct {
	Image  string `json:"image"`
	Status string `json:"status"` // "pulled" or "up-to-date"
}

// MysqlImagePullOutput is the output schema for the mysql_image_pull tool.
type MysqlImagePullOutput struct {
	Pulled []MysqlImagePullEntry `json:"pulled"`
}

// NewMysqlImagePull returns a tool handler for mysql_image_pull using the given Docker client ID.
func NewMysqlImagePull(clientId docker.DockerClientId) func(context.Context, *mcp.CallToolRequest, MysqlImagePullInput) (*mcp.CallToolResult, MysqlImagePullOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, in MysqlImagePullInput) (*mcp.CallToolResult, MysqlImagePullOutput, error) {
		ctx, skprClient, err := client.New(ctx)
		if err != nil {
			return nil, MysqlImagePullOutput{}, fmt.Errorf("failed to create client: %w", err)
		}

		getRepositoryResp, err := skprClient.Mysql().ImageGetRepository(ctx, &pb.ImageGetRepositoryRequest{
			Environment: in.Environment,
		})
		if err != nil {
			return nil, MysqlImagePullOutput{}, fmt.Errorf("failed to get repository: %w", err)
		}

		auth := types.Auth{
			Username: skprClient.Credentials.Username,
			Password: skprClient.Credentials.Password,
			Session:  skprClient.Credentials.Session,
		}

		if ecr.IsRegistry(getRepositoryResp.Repository) {
			auth, err = ecr.UpgradeAuth(ctx, getRepositoryResp.Repository, auth)
			if err != nil {
				return nil, MysqlImagePullOutput{}, errors.Wrap(err, "failed to upgrade AWS ECR authentication")
			}
		}

		dockerClient, err := docker.NewClientFromUserConfig(auth, clientId)
		if err != nil {
			return nil, MysqlImagePullOutput{}, errors.Wrap(err, "failed to create Docker client")
		}

		// Build tag list — same logic as the CLI command.
		tags := []string{}
		if in.ID != "" {
			tags = append(tags, in.ID)
		} else {
			databases := in.Databases
			if len(databases) == 0 {
				databases = []string{"default"}
			}
			for _, database := range databases {
				tags = append(tags, fmt.Sprintf("%s-latest", database))
			}
		}

		// MCP stdout must only carry JSON-RPC frames; route Docker progress to
		// stderr so it doesn't corrupt the transport.
		progressWriter := io.Writer(os.Stderr)

		out := MysqlImagePullOutput{
			Pulled: make([]MysqlImagePullEntry, 0, len(tags)),
		}

		for _, tag := range tags {
			imageName := fmt.Sprintf("%s:%s", getRepositoryResp.Repository, tag)

			cleanupId, err := dockerClient.ImageId(context.TODO(), imageName)
			if err != nil {
				return nil, MysqlImagePullOutput{}, err
			}

			err = dockerClient.PullImage(context.TODO(), getRepositoryResp.Repository, tag, progressWriter)
			if err != nil {
				return nil, MysqlImagePullOutput{}, err
			}

			currentId, err := dockerClient.ImageId(context.TODO(), imageName)
			if err != nil {
				return nil, MysqlImagePullOutput{}, err
			}

			status := "pulled"
			if cleanupId == currentId {
				status = "up-to-date"
			}

			out.Pulled = append(out.Pulled, MysqlImagePullEntry{
				Image:  imageName,
				Status: status,
			})

			// Clean up the old image if it differs from the new one.
			if cleanupId != "" && cleanupId != currentId {
				if err := dockerClient.RemoveImage(context.TODO(), cleanupId); err != nil {
					fmt.Fprintf(os.Stderr, "warning: failed to remove old image %s: %v\n", cleanupId, err)
				}
			}
		}

		// Summarise for the text part of the tool result, keeping stdout JSON-RPC-clean.
		summary := make([]string, 0, len(out.Pulled))
		for _, entry := range out.Pulled {
			summary = append(summary, fmt.Sprintf("%s (%s)", entry.Image, entry.Status))
		}

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: strings.Join(summary, "\n")},
			},
		}

		return result, out, nil
	}
}
