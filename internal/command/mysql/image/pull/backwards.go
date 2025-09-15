package pull

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gosuri/uilive"
	"github.com/pkg/errors"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/internal/client"
)

// @todo, Remove in a future release.
func (cmd *Command) runBackwardsCompat() error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	project, err := client.Project().Get(ctx, &pb.ProjectGetRequest{})
	if err != nil {
		return errors.Wrap(err, "failed to get project")
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
	if ecr.IsRegistry(project.Registry.MySQL) {
		auth, err = ecr.UpgradeAuth(ctx, project.Registry.MySQL, creds)
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

	var images []string

	for _, database := range cmd.Params.Databases {
		tag, err := getTag(ctx, client.Mysql(), database, cmd.Params)
		if err != nil {
			return fmt.Errorf("failed to compute image tag: %w", err)
		}

		imageName := fmt.Sprintf("%s:%s", project.Registry.MySQL, tag)

		// Keep for later so we can inform the developer on which images they can use.
		images = append(images, imageName)

		fmt.Printf("Pulling: %s\n", imageName)

		// Lookup the ID of the current image so we can delete it after we pull the image one.
		cleanup, err := dockerclient.InspectImage(imageName)
		if err != nil && !errors.Is(err, docker.ErrNoSuchImage) {
			return err
		}

		opts := docker.PullImageOptions{
			OutputStream: writer,
			Repository:   project.Registry.MySQL,
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

		fmt.Println("Cleaning up old image with:", cleanup.ID)

		err = dockerclient.RemoveImage(cleanup.ID)
		if err != nil {
			return err
		}
	}

	writer.Stop()

	fmt.Println("The following images have been successfully pulled:")

	for _, image := range images {
		fmt.Println("-", image)
	}

	return nil
}
