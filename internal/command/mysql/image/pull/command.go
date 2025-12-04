package pull

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gosuri/uilive"
	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/auth"
	"github.com/skpr/cli/internal/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/docker"
	"github.com/skpr/cli/internal/docker/dockerclient"
	"github.com/skpr/cli/internal/docker/goclient"
	skprlog "github.com/skpr/cli/internal/log"
)

const (
	// DefaultTagSuffix used to pull the latest database image built by the system.
	DefaultTagSuffix = "latest"
)

// Command to pull a database image.
type Command struct {
	Params Params
}

// Params provided to this command.
type Params struct {
	Environment string
	Databases   []string
	ID          string
	Tag         string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	prettyHandler := skprlog.NewHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   false,
		ReplaceAttr: nil,
	})

	logger := slog.New(prettyHandler)

	getRepositoryResp, err := client.Mysql().ImageGetRepository(ctx, &pb.ImageGetRepositoryRequest{
		Environment: cmd.Params.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	auth := auth.Auth{
		Username: client.Credentials.Username,
		Password: client.Credentials.Password,
	}

	// @todo, Consider abstracting this if another registry + credentials pair is required.
	if ecr.IsRegistry(getRepositoryResp.Repository) {
		auth, err = ecr.UpgradeAuth(ctx, getRepositoryResp.Repository, client.Credentials)
		if err != nil {
			return errors.Wrap(err, "failed to upgrade AWS ECR authentication")
		}
	}

	c, err := getClient(auth)
	if err != nil {
		return errors.Wrap(err, "failed to create Docker client")
	}

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for _, database := range cmd.Params.Databases {
		tag := fmt.Sprintf("%s-%s", database, DefaultTagSuffix)

		imageName := fmt.Sprintf("%s:%s", getRepositoryResp.Repository, tag)

		logger.Info(fmt.Sprintf("Pulling: %s", imageName))

		// Lookup the ID of the current image so we can delete it after we pull the image one.
		cleanupId, err := c.ImageId(context.TODO(), imageName)
		if err != nil {
			return err
		}

		err = c.PullImage(context.TODO(), getRepositoryResp.Repository, tag, writer)
		if err != nil {
			return err
		}

		// Check if there was an old image before cleaning up.
		if cleanupId == "" {
			continue
		}

		currentId, err := c.ImageId(context.TODO(), imageName)
		if err != nil {
			return err
		}

		// Don't cleanup the old image if it was the latest and never needed to be updated.
		if cleanupId == currentId {
			continue
		}

		logger.Info(fmt.Sprintf("Cleaning up old image with: %s", cleanupId))

		err = c.RemoveImage(context.TODO(), cleanupId)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Successfully pulled image: %s:%s", getRepositoryResp.Repository, tag))
	}

	return nil
}

func getClient(auth auth.Auth) (docker.DockerClient, error) {
	// See if we're using default builder.
	userConfig, _ := user.NewClient()
	featureFlags, _ := userConfig.LoadFeatureFlags()

	if featureFlags.Builder == user.ConfigPackageBuilderDocker {
		c, err := dockerclient.New(auth)
		return c, err
	}

	if featureFlags.Builder != "" && featureFlags.Builder != user.ConfigPackageBuilderLegacy {
		return nil, fmt.Errorf("unknown docker client: %s", featureFlags.Builder)
	}

	return goclient.New(auth)
}
