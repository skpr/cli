package buildpack

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/egym-playground/go-prefix-writer/prefixer"
	docker "github.com/fsouza/go-dockerclient"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/buildpack/utils/image"
	"github.com/skpr/cli/internal/color"
)

// DockerClientInterface provides an interface that allows us to test the builder.
type DockerClientInterface interface {
	BuildImage(options docker.BuildImageOptions) error
	PushImage(options docker.PushImageOptions, auth docker.AuthConfiguration) error
	InspectImage(name string) (*docker.Image, error)
}

// Builder is the docker image builder.
type Builder struct {
	dockerClient DockerClientInterface
}

// BuildResponse is returned by the build operation.
type BuildResponse struct {
	Images []Image `json:"images"`
}

// ImageType used to identify what a built image is used for.
type ImageType string

const (
	// ImageTypeCompile is used to identify images which were built during the "compile" phase.
	ImageTypeCompile ImageType = "compile"
	// ImageTypeRuntime is used to identify images which were built for "runtime".
	ImageTypeRuntime ImageType = "runtime"
)

// Image build has been built.
type Image struct {
	// Name of the image.
	Name string `json:"name"`
	// Type of image that has been built.
	Type ImageType `json:"type"`
	// Tag used to push image.
	Tag string `json:"tag"`
	// Digest of the image.
	Digest string `json:"digest"`
}

// Params used for building the applications.
type Params struct {
	Auth     docker.AuthConfiguration
	Writer   io.Writer
	Context  string
	Registry string
	NoPush   bool
	Version  string
}

// Dockerfiles the docker build files.
type Dockerfiles map[string]string

const (
	// ImageNameCompile is used for compiling the application.
	ImageNameCompile = "compile"

	// BuildArgCompileImage is used for referencing the compile image.
	BuildArgCompileImage = "COMPILE_IMAGE"
	// BuildArgVersion is used for providing the version identifier of the application.
	BuildArgVersion = "SKPR_VERSION"
)

// NewBuilder creates a new Builder.
func NewBuilder(dockerClient DockerClientInterface) *Builder {
	return &Builder{
		dockerClient: dockerClient,
	}
}

// Build the images.
func (b *Builder) Build(dockerfiles Dockerfiles, params Params) (BuildResponse, error) {
	var resp BuildResponse

	compileDockerfile, ok := dockerfiles[ImageNameCompile]
	if !ok {
		return resp, fmt.Errorf("%q is a required dockerfile", ImageNameCompile)
	}

	args := []docker.BuildArg{
		{
			Name:  BuildArgVersion,
			Value: params.Version,
		},
	}

	start := time.Now()

	// We build the compile image first, as it is the base image for other images.
	compileBuild := docker.BuildImageOptions{
		Name:         image.Name(params.Registry, params.Version, ImageNameCompile),
		Dockerfile:   compileDockerfile,
		ContextDir:   params.Context,
		OutputStream: prefixWithTime(params.Writer, ImageNameCompile, start),
		BuildArgs:    args,
	}

	resp.Images = append(resp.Images, Image{
		Name: ImageNameCompile,
		Type: ImageTypeCompile,
		Tag:  compileBuild.Name,
	})

	// We need to build the 'compile' image first.
	fmt.Fprintf(params.Writer, "Building image: %s\n", compileBuild.Name)
	err := b.dockerClient.BuildImage(compileBuild)
	if err != nil {
		return resp, err
	}
	fmt.Fprintf(params.Writer, "Built compile image in %s\n", time.Since(start).Round(time.Second))

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, ImageNameCompile)

	// Adds compile image identifier to the runtime images as an arg.
	// That allows runtime images to copy over the compiled code.
	args = append(args, docker.BuildArg{
		Name:  BuildArgCompileImage,
		Value: image.Name(params.Registry, params.Version, ImageNameCompile),
	})

	var builds []docker.BuildImageOptions

	for imageName, dockerfile := range dockerfiles {
		build := docker.BuildImageOptions{
			Name:         image.Name(params.Registry, params.Version, imageName),
			Dockerfile:   dockerfile,
			ContextDir:   params.Context,
			OutputStream: prefixWithTime(params.Writer, imageName, start),
			BuildArgs:    args,
		}

		// Add to the builder list.
		builds = append(builds, build)

		// Add to the manifest.
		resp.Images = append(resp.Images, Image{
			Name: imageName,
			Type: ImageTypeRuntime,
			Tag:  build.Name,
		})
	}

	bg, ctx := errgroup.WithContext(context.Background())

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

	for imageName := range dockerfiles {
		// Compile image is only for building, so we don't push.
		if imageName == ImageNameCompile {
			continue
		}

		pushes = append(pushes, docker.PushImageOptions{
			Name: params.Registry,
			Tag:  image.Tag(params.Version, imageName),
		})
	}

	pg, ctx := errgroup.WithContext(context.Background())

	for _, push := range pushes {
		// https://golang.org/doc/faq#closures_and_goroutines
		push := push

		// Allows us to cancel push executions.
		push.Context = ctx

		fmt.Fprintf(params.Writer, "Pushing image: %s:%s\n", push.Name, push.Tag)

		pg.Go(func() error {
			start = time.Now()

			err = b.dockerClient.PushImage(push, params.Auth)
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

	var images []Image

	for _, respImage := range resp.Images {
		// Compile image is only for building, so we don't push.
		if respImage.Name == ImageNameCompile {
			continue
		}

		fmt.Fprintf(params.Writer, "Fetching digest for: %s\n", respImage.Name)

		inspect, err := b.dockerClient.InspectImage(image.Name(params.Registry, params.Version, respImage.Name))
		if err != nil {
			return resp, fmt.Errorf("failed to inspect image %q: %w", respImage.Name, err)
		}

		digest, err := getDigest(inspect.RepoDigests)
		if err != nil {
			return resp, fmt.Errorf("failed to get digest for image %q: %w", respImage.Name, err)
		}

		respImage.Digest = digest

		images = append(images, respImage)
	}

	resp.Images = images

	fmt.Fprintf(params.Writer, "Build complete in: %s\n", time.Since(start).Round(time.Second))

	return resp, nil
}

// Helper function to prefix all output for a stream.
func prefixWithTime(w io.Writer, name string, start time.Time) io.Writer {
	return prefixer.New(w, newPrefixer(color.Wrap(strings.ToUpper(name)), start).PrefixFunc())
}

// Return a digest from a list of digests.
func getDigest(digests []string) (string, error) {
	if len(digests) == 0 {
		return "", fmt.Errorf("digest not found")
	}

	// Take the first one off the list.
	// https://notaryproject.dev/docs/quickstart-guides/quickstart-sign-image-artifact/#add-an-image-to-the-oci-compatible-registry
	sl := strings.Split(digests[0], "@")

	if len(sl) != 2 {
		return "", fmt.Errorf("invalid digest format")
	}

	return sl[1], nil
}
