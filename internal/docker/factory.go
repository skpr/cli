package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/skpr/cli/internal/auth"
	"github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/docker/dockerclient"
	"github.com/skpr/cli/internal/docker/goclient"
	"github.com/skpr/cli/internal/docker/mockclient"
)

type DockerClient interface {
	ImageId(ctx context.Context, name string) (string, error)
	PullImage(ctx context.Context, repository, tag string, writer io.Writer) error
	PushImage(ctx context.Context, repository, tag string, writer io.Writer) error
	RemoveImage(ctx context.Context, id string) error
	BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error
}

func NewClientFromUserConfig(auth auth.Auth) (DockerClient, error) {
	// See if we're using default builder.
	userConfig, _ := user.NewClient()
	featureFlags, _ := userConfig.LoadFeatureFlags()

	switch featureFlags.DockerClient {
	case user.ConfigPackageClientLegacy:
		return goclient.New(auth)
	case user.ConfigPackageClientDocker:
		return dockerclient.New(auth)
	case user.ConfigPackageClientMock:
		return mockclient.New(), nil
	}

	if featureFlags.DockerClient != "" {
		return nil, fmt.Errorf("unknown docker client: %s", featureFlags.DockerClient)
	}

	return goclient.New(auth)
}
