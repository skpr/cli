package pull

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gosuri/uilive"
	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/internal/client"
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
		logger.Info("Using backwards compatibility command to pull image.")
		return cmd.runBackwardsCompat()
	}

	creds, err := client.CredentialsProvider.Retrieve(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}

	auth := docker.AuthConfiguration{
		Username: creds.AccessKeyID,
		Password: creds.SecretAccessKey,
	}

	// @todo, Consider abstracting this if another registry + credentials pair is required.
	if ecr.IsRegistry(getRepositoryResp.Repository) {
		auth, err = ecr.UpgradeAuth(ctx, getRepositoryResp.Repository, creds)
		if err != nil {
			return errors.Wrap(err, "failed to upgrade AWS ECR authentication")
		}
	}

	dockerclient, err := docker.NewClientFromEnv()
	if err != nil {
		return errors.Wrap(err, "failed to setup Docker client")
	}

	writer := uilive.New()
	writer.Start()

	for _, database := range cmd.Params.Databases {
		tag := fmt.Sprintf("%s-%s", database, DefaultTagSuffix)

		imageName := fmt.Sprintf("%s:%s", getRepositoryResp.Repository, tag)

		logger.Info(fmt.Sprintf("Pulling: %s", imageName))

		// Lookup the ID of the current image so we can delete it after we pull the image one.
		cleanup, err := dockerclient.InspectImage(imageName)
		if err != nil && !errors.Is(err, docker.ErrNoSuchImage) {
			return err
		}

		opts := docker.PullImageOptions{
			OutputStream: writer,
			Repository:   getRepositoryResp.Repository,
			Tag:          tag,
		}

		err = dockerclient.PullImage(opts, auth)
		if err != nil {
			return err
		}

		// Check if there was an old images which
		if cleanup == nil {
			continue
		}

		current, err := dockerclient.InspectImage(imageName)
		if err != nil {
			return err
		}

		// Don't cleanup the old image if it was the latest and never needed to be updated.
		if cleanup.ID == current.ID {
			continue
		}

		logger.Info(fmt.Sprintf("Cleaning up old image with: %s", cleanup.ID))

		err = dockerclient.RemoveImage(cleanup.ID)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Successfully pulled image: %s:%s", getRepositoryResp.Repository, tag))
	}

	writer.Stop()

	return nil
}
