package buildpack

import (
	"context"
	"fmt"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/buildpack/types"
	"github.com/skpr/cli/internal/buildpack/utils/finder"
	"github.com/skpr/cli/internal/buildpack/utils/image"
	"github.com/skpr/cli/internal/buildpack/utils/prefixer"
)

// DockerClientInterface provides an interface that allows us to test the builder.
type DockerClientInterface interface {
	BuildImage(options docker.BuildImageOptions) error
	PushImage(options docker.PushImageOptions, auth docker.AuthConfiguration) error
}

// Builder is the docker image builder.
type Builder struct {
	dockerClient DockerClientInterface
}

// NewBuilder creates a new Builder.
func NewBuilder() (*Builder, error) {
	dockerclient, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to setup Docker client: %w", err)
	}

	return &Builder{
		dockerClient: dockerclient,
	}, nil
}

// Build the images.
func (b *Builder) Build(ctx context.Context, dockerfiles finder.Dockerfiles, params types.Params) (types.BuildResponse, error) {
	var resp types.BuildResponse

	compileDockerfile, ok := dockerfiles[types.ImageNameCompile]
	if !ok {
		return resp, fmt.Errorf("%q is a required dockerfile", types.ImageNameCompile)
	}

	args := []docker.BuildArg{
		{
			Name:  types.BuildArgVersion,
			Value: params.Version,
		},
	}

	for k, v := range params.BuildArgs {
		args = append(args, docker.BuildArg{
			Name:  k,
			Value: v,
		})
	}

	start := time.Now()

	// We build the compile image first, as it is the base image for other images.
	compileBuild := docker.BuildImageOptions{
		Name:         image.Name(params.Registry, params.Version, types.ImageNameCompile),
		Dockerfile:   compileDockerfile,
		ContextDir:   params.Context,
		OutputStream: prefixer.WrapWriterWithPrefixer(params.Writer, types.ImageNameCompile, start),
		BuildArgs:    args,
	}

	// We need to build the 'compile' image first.
	fmt.Fprintf(params.Writer, "Building image: %s\n", compileBuild.Name)
	err := b.dockerClient.BuildImage(compileBuild)
	if err != nil {
		return resp, err
	}
	fmt.Fprintf(params.Writer, "Built compile image in %s\n", time.Since(start).Round(time.Second))

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, types.ImageNameCompile)

	// Adds compile image identifier to the runtime images as an arg.
	// That allows runtime images to copy over the compiled code.
	args = append(args, docker.BuildArg{
		Name:  types.BuildArgCompileImage,
		Value: image.Name(params.Registry, params.Version, types.ImageNameCompile),
	})

	var builds []docker.BuildImageOptions

	for imageName, dockerfile := range dockerfiles {
		build := docker.BuildImageOptions{
			Name:         image.Name(params.Registry, params.Version, imageName),
			Dockerfile:   dockerfile,
			ContextDir:   params.Context,
			OutputStream: prefixer.WrapWriterWithPrefixer(params.Writer, imageName, start),
			BuildArgs:    args,
		}

		// Add to the builder list.
		builds = append(builds, build)

		// Add to the manifest.
		resp.Images = append(resp.Images, types.Image{
			Name: imageName,
			Tag:  build.Name,
		})
	}

	bg, ctx := errgroup.WithContext(ctx)

	for _, build := range builds {
		// https://golang.org/doc/faq#closures_and_goroutines
		build := build

		// Allows us to cancel build executions.
		build.Context = ctx

		fmt.Fprintf(params.Writer, "Building image: %s\n", build.Name)

		bg.Go(func() error {
			start = time.Now()

			err := b.dockerClient.BuildImage(build)
			if err != nil {
				return err
			}

			fmt.Fprintf(params.Writer, "Built %s image in %s\n", build.Name, time.Since(start).Round(time.Second))

			return nil
		})
	}

	err = bg.Wait()
	if err != nil {
		return resp, err
	}

	if params.NoPush {
		return resp, nil
	}

	var pushes []docker.PushImageOptions

	for _, buildImage := range resp.Images {
		pushes = append(pushes, docker.PushImageOptions{
			Name: params.Registry,
			Tag:  image.Tag(params.Version, buildImage.Name),
		})
	}

	pg, ctx := errgroup.WithContext(ctx)

	auth := docker.AuthConfiguration{
		Username: params.Auth.Username,
		Password: params.Auth.Password,
	}

	for _, push := range pushes {
		// https://golang.org/doc/faq#closures_and_goroutines
		push := push

		// Allows us to cancel push executions.
		push.Context = ctx

		fmt.Fprintf(params.Writer, "Pushing image: %s:%s\n", push.Name, push.Tag)

		pg.Go(func() error {
			start = time.Now()

			err = b.dockerClient.PushImage(push, auth)
			if err != nil {
				return err
			}

			fmt.Fprintf(params.Writer, "Pushed %s:%s image in %s\n", push.Name, push.Tag, time.Since(start).Round(time.Second))

			return nil
		})
	}

	err = pg.Wait()
	if err != nil {
		return resp, err
	}

	fmt.Fprintf(params.Writer, "Build complete in: %s\n", time.Since(start).Round(time.Second))

	return resp, nil
}
