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
	"github.com/skpr/cli/internal/docker"
	skprlog "github.com/skpr/cli/internal/log"
)

const (
	// DefaultTagSuffix used to pull the latest database image built by the system.
	DefaultTagSuffix = "latest"
)

// Command to pull a database image.
type Command struct {
	Params   Params
	ClientId docker.DockerClientId
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

	c, err := docker.NewClientFromUserConfig(auth, cmd.ClientId)
	if err != nil {
		return errors.Wrap(err, "failed to create Docker client")
	}

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	tags := []string{}
	if cmd.Params.ID != "" {
		tags = append(tags, cmd.Params.ID)
	} else {
		for _, database := range cmd.Params.Databases {
			tag := fmt.Sprintf("%s-%s", database, DefaultTagSuffix)
			tags = append(tags, tag)
		}
	}

	for _, tag := range tags {
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

		currentId, err := c.ImageId(context.TODO(), imageName)
		if err != nil {
			return err
		}

		if cleanupId == currentId {
			logger.Info(fmt.Sprintf("Image is up to date: %s", imageName))
		} else {
			logger.Info(fmt.Sprintf("Successfully pulled image: %s", imageName))
		}

		// If it's a fresh image or the same image as the current one, skip deleting it.
		if cleanupId == "" || cleanupId == currentId {
			continue
		}

		logger.Info(fmt.Sprintf("Cleaning up old image with: %s", cleanupId))

		err = c.RemoveImage(context.TODO(), cleanupId)
		if err != nil {
			return err
		}
	}

	return nil
}
