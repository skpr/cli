package docker

import (
	"context"
	"io"
)

type DockerClient interface {
	ImageId(ctx context.Context, name string) (string, error)
	PullImage(ctx context.Context, repository, tag string, writer io.Writer) error
	PushImage(ctx context.Context, repository, tag string, writer io.Writer) error
	RemoveImage(ctx context.Context, id string) error
	BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error
}
