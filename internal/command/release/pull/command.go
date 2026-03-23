package pull

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gosuri/uilive"
	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/containers/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/containers/docker"
	"github.com/skpr/cli/containers/docker/types"
	"github.com/skpr/cli/internal/client"
	skprlog "github.com/skpr/cli/internal/log"
)

// Command to pull a database image.
type Command struct {
	Params   Params
	ClientId docker.DockerClientId
}

// Params provided to this command.
type Params struct {
	Name string
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

	release, err := client.Release().Info(ctx, &pb.ReleaseInfoRequest{
		Name: cmd.Params.Name,
	})
	if err != nil {
		return fmt.Errorf("could not get release: %w", err)
	}

	for _, image := range release.Images {
		repository, tag, err := ParseImage(image.URI)
		if err != nil {
			return errors.Wrap(err, "failed to parse image reference")
		}

		auth := types.Auth{
			Username: client.Credentials.Username,
			Password: client.Credentials.Password,
			Session:  client.Credentials.Session,
		}

		// @todo, Consider abstracting this if another registry + credentials pair is required.
		if ecr.IsRegistry(repository) {
			auth, err = ecr.UpgradeAuth(ctx, repository, auth)
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

		logger.Info(fmt.Sprintf("Pulling: %s", image.URI))

		err = c.PullImage(ctx, repository, tag, writer)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Successfully pulled image: %s", image.URI))
	}

	return nil
}

func ParseImage(image string) (repository string, tag string, err error) {
	if image == "" {
		return "", "", fmt.Errorf("image reference is empty")
	}

	// Reject digest references explicitly
	if strings.Contains(image, "@") {
		return "", "", fmt.Errorf("digest references are not supported")
	}

	// Split on the last colon to preserve registry ports
	lastColon := strings.LastIndex(image, ":")
	if lastColon == -1 {
		return "", "", fmt.Errorf("image reference does not contain a tag")
	}

	repository = image[:lastColon]
	tag = image[lastColon+1:]

	if repository == "" {
		return "", "", fmt.Errorf("repository is empty")
	}
	if tag == "" {
		return "", "", fmt.Errorf("tag is empty")
	}

	return repository, tag, nil
}
