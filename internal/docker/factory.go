package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/skpr/cli/internal/auth"
	"github.com/skpr/cli/internal/docker/dockerclient"
	"github.com/skpr/cli/internal/docker/goclient"
	"github.com/skpr/cli/internal/docker/mockclient"
)

type DockerClient interface {
	ImageId(ctx context.Context, name string) (string, error)
	PullImage(ctx context.Context, registry, tag string, writer io.Writer) error
	PushImage(ctx context.Context, registry, tag string, writer io.Writer) error
	RemoveImage(ctx context.Context, id string) error
	BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error
}

type DockerClientId string

const (
	DockerClientIdLegacy DockerClientId = "legacy"
	DockerClientIdDocker DockerClientId = "docker"
	DockerClientIdMock   DockerClientId = "mock"
)

func NewClientFromUserConfig(auth auth.Auth, clientId DockerClientId) (DockerClient, error) {
	switch clientId {
	case DockerClientIdLegacy:
		return goclient.New(auth)
	case DockerClientIdDocker:
		return dockerclient.New(auth)
	case DockerClientIdMock:
		return mockclient.New(), nil
	}

	if clientId != "" {
		return nil, fmt.Errorf("unknown docker client: %s", clientId)
	}

	return dockerclient.New(auth)
}
