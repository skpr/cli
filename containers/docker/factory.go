package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/skpr/cli/containers/docker/dockerclient"
	"github.com/skpr/cli/containers/docker/mockclient"
	"github.com/skpr/cli/containers/docker/types"
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
	DockerClientIdDocker DockerClientId = "docker"
	DockerClientIdMock   DockerClientId = "mock"
)

func NewClientFromUserConfig(auth types.Auth, clientId DockerClientId) (DockerClient, error) {
	switch clientId {
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
