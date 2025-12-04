package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/skpr/api/pb"

	buildpack "github.com/skpr/cli/internal/buildpack/builder"
	"github.com/skpr/cli/internal/buildpack/utils/aws/ecr"
	"github.com/skpr/cli/internal/buildpack/utils/finder"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config"
	"github.com/skpr/cli/internal/slice"
)

// Command to package an application.
type Command struct {
	Region        string
	PackageDir    string
	Params        buildpack.Params
	PrintManifest bool
	BuildArgs     []string
	Platform      string
	Debug         bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	if cmd.Params.NoPush {
		config, err := config.New()
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}

		cmd.Params.Registry = fmt.Sprintf("localhost/skpr/%s", config.Project)
	} else {
		ctx, client, err := client.New(ctx)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		project, err := client.Project().Get(ctx, &pb.ProjectGetRequest{})
		if err != nil {
			return fmt.Errorf("failed to list environments: %w", err)
		}

		releases, err := client.Release().List(ctx, &pb.ReleaseListRequest{})
		if err != nil {
			return fmt.Errorf("failed to list releases: %w", err)
		}

		for _, release := range releases.Items {
			if release.Name == cmd.Params.Version {
				return fmt.Errorf("release has already been created")
			}
		}

		cmd.Params.Registry = project.Registry.Application

		cmd.Params.Auth = docker.AuthConfiguration{
			Username: client.Credentials.Username,
			Password: client.Credentials.Password,
		}

		// @todo, Consider abstracting this if another registry + credentials pair is required.
		if ecr.IsRegistry(cmd.Params.Registry) {
			auth, err := ecr.UpgradeAuth(ctx, cmd.Params.Registry, client.Credentials)
			if err != nil {
				return fmt.Errorf("failed to upgrade AWS ECR authentication: %w", err)
			}

			cmd.Params.Auth = auth
		}
	}

	cmd.Params.Writer = os.Stderr
	cmd.Params.Platform = cmd.Platform

	// Convert build args from slice to map.
	//   eg. --build-arg=KEY=VALUE to map[string]string{"KEY": "VALUE"}
	cmd.Params.BuildArgs = slice.ToMap(cmd.BuildArgs, "=")

	dockerfiles, err := finder.FindDockerfiles(cmd.PackageDir)
	if err != nil {
		return fmt.Errorf("failed to find dockerfiles: %w", err)
	}

	if cmd.Debug {
		fmt.Println("Found the following dockerfiles:")
		for key, path := range dockerfiles {
			fmt.Printf("%-10s %q\n", key, path)
		}
	}

	// Print deprecation notice.
	for key, path := range dockerfiles {
		if strings.HasSuffix(path, ".dockerfile") {
			fmt.Printf("[DEPRECATED] Dockerfile location %q is deprecated. Use \"%s/%s/Dockerfile\" instead.\n", path, filepath.Dir(path), key)
		}
	}

	dockerclient, err := docker.NewClientFromEnv()
	if err != nil {
		return fmt.Errorf("failed to setup Docker client: %w", err)
	}

	builder := buildpack.NewBuilder(dockerclient)

	resp, err := builder.Build(dockerfiles, cmd.Params)
	if err != nil {
		return err
	}

	request := &pb.ReleaseCreateRequest{
		Name: cmd.Params.Version,
	}

	if cmd.PrintManifest {
		b, err := json.MarshalIndent(resp.Images, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(b))
	}

	// Do not create a release if the --no-push flag is set.
	if cmd.Params.NoPush {
		return nil
	}

	request = &pb.ReleaseCreateRequest{
		Name: cmd.Params.Version,
	}

	for _, image := range resp.Images {
		request.Images = append(request.Images, &pb.ReleaseImage{
			Name: image.Name,
			URI:  image.Tag,
		})
	}

	ctx, client, err := client.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	_, err = client.Release().Create(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	if !cmd.PrintManifest {
		fmt.Printf("Release %v was successfully created, and can be deployed using `skpr deploy <env> %v`.\n", cmd.Params.Version, cmd.Params.Version)
	}

	return nil
}
