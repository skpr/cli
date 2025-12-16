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
	ClientIdLegacy DockerClientId = "legacy"
	ClientIdDocker DockerClientId = "docker"
	ClientIdMock   DockerClientId = "mock"
)

func NewClientFromUserConfig(auth auth.Auth, clientId DockerClientId) (DockerClient, error) {
	switch clientId {
	case ClientIdLegacy:
		return goclient.New(auth)
	case ClientIdDocker:
		return dockerclient.New(auth)
	case ClientIdMock:
		return mockclient.New(), nil
	}

	if clientId != "" {
		return nil, fmt.Errorf("unknown docker client: %s", clientId)
	}

	return goclient.New(auth)
}
