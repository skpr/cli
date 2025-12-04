package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/skpr/cli/internal/auth"
	"github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/docker/dockerclient"
	"github.com/skpr/cli/internal/docker/goclient"
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

	if featureFlags.DockerClient == user.ConfigPackageClientDocker {
		c, err := dockerclient.New(auth)
		return c, err
	} else if featureFlags.DockerClient == user.ConfigPackageClientMock {
		c, err := dockerclient.New(auth)
		return c, err
	}

	if featureFlags.DockerClient != "" && featureFlags.DockerClient != user.ConfigPackageClientLegacy {
		return nil, fmt.Errorf("unknown docker client: %s", featureFlags.DockerClient)
	}

	return goclient.New(auth)
}
