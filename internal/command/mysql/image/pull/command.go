package pull

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	imagetypes "github.com/docker/docker/api/types/image"
	dockregistry "github.com/docker/docker/api/types/registry"
	dockclient "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
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
		return fmt.Errorf("failed to get repository: %w", err)
	}

	auth := dockregistry.AuthConfig{
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

	// Official Docker SDK client
	dc, err := dockclient.NewClientWithOpts(dockclient.FromEnv, dockclient.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "failed to setup Docker client")
	}

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	authHdr, err := encodeRegistryAuth(auth)
	if err != nil {
		return errors.Wrap(err, "failed to encode registry auth")
	}

	for _, database := range cmd.Params.Databases {
		tag := fmt.Sprintf("%s-%s", database, DefaultTagSuffix)
		imageName := fmt.Sprintf("%s:%s", getRepositoryResp.Repository, tag)

		logger.Info(fmt.Sprintf("Pulling: %s", imageName))

		// Lookup the ID of the current (pre-pull) image so we can delete it after we pull the new one.
		var oldID string
		if insp, _, err := dc.ImageInspectWithRaw(ctx, imageName); err == nil {
			oldID = insp.ID
		} else if !errdefs.IsNotFound(err) {
			return err
		}

		rc, err := dc.ImagePull(ctx, imageName, imagetypes.PullOptions{
			RegistryAuth: authHdr,
		})
		if err != nil {
			return err
		}
		// Stream the pull progress to the UI writer.
		if _, copyErr := io.Copy(writer, rc); copyErr != nil {
			_ = rc.Close()
			return copyErr
		}
		_ = rc.Close()

		// If there was no old image, we’re done.
		if oldID == "" {
			logger.Info(fmt.Sprintf("Successfully pulled image: %s", imageName))
			continue
		}

		// Inspect current image to compare IDs.
		cur, _, err := dc.ImageInspectWithRaw(ctx, imageName)
		if err != nil {
			return err
		}

		// Don't cleanup the old image if it was unchanged.
		if oldID == cur.ID {
			logger.Info(fmt.Sprintf("Image already up to date: %s", imageName))
			continue
		}

		logger.Info(fmt.Sprintf("Cleaning up old image with ID: %s", oldID))
		_, err = dc.ImageRemove(ctx, oldID, imagetypes.RemoveOptions{
			PruneChildren: true,
			Force:         false,
		})
		if err != nil && !errdefs.IsNotFound(err) {
			return err
		}

		logger.Info(fmt.Sprintf("Successfully pulled image: %s", imageName))
	}

	return nil
}

// encodeRegistryAuth converts a registry.AuthConfig to the base64 JSON string
// expected by the Docker Engine API for pull/push operations.
func encodeRegistryAuth(cfg dockregistry.AuthConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
