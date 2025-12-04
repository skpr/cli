package buildpack

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/moby/patternmatcher/ignorefile"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/buildpack/types"
	"github.com/skpr/cli/internal/buildpack/utils/finder"
	"github.com/skpr/cli/internal/buildpack/utils/image"
	"github.com/skpr/cli/internal/buildpack/utils/prefixer"
	"github.com/skpr/cli/internal/docker"
)

// Builder is the docker image builder.
type Builder struct {
	Client docker.DockerClient
}

// NewBuilder creates a new Builder.
func NewBuilder(c docker.DockerClient) (*Builder, error) {
	return &Builder{
		Client: c,
	}, nil
}

// Build the images.
func (b *Builder) Build(ctx context.Context, dockerfiles finder.Dockerfiles, params types.Params) (types.BuildResponse, error) {
	var resp types.BuildResponse

	excludePatterns, err := loadIgnoreFilePatterns(params.IgnoreFile)
	if err != nil {
		return resp, fmt.Errorf("failed to parse ignore file: %w", err)
	}

	compileDockerfile, ok := dockerfiles[types.ImageNameCompile]
	if !ok {
		return resp, fmt.Errorf("%q is a required dockerfile", types.ImageNameCompile)
	}

	// Build the build args.
	buildArgs := params.BuildArgs
	buildArgs[types.BuildArgVersion] = params.Version

	start := time.Now()

	// Build the compile image first; it's the base for others.
	compileRef := image.Name(params.Registry, params.Version, types.ImageNameCompile)

	fmt.Fprintf(params.Writer, "Building image: %s\n", compileRef)

	localOut := prefixer.WrapWriterWithPrefixer(params.Writer, types.ImageNameCompile, start)
	err = b.Client.BuildImage(ctx, compileDockerfile, params.Context, compileRef, excludePatterns, buildArgs, localOut)
	if err != nil {
		return resp, err
	}

	fmt.Fprintf(params.Writer, "Built %s image in %s\n", compileRef, time.Since(start).Round(time.Second))

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, types.ImageNameCompile)

	// Adds compile image identifier to the runtime images as an arg.
	buildArgs[types.BuildArgCompileImage] = compileRef

	// Prepare runtime builds.
	type pendingBuild struct {
		name       string
		imageRef   string
		dockerfile string
	}
	var builds []pendingBuild
	for imageName, dockerfile := range dockerfiles {
		ref := image.Name(params.Registry, params.Version, imageName)
		builds = append(builds, pendingBuild{name: imageName, imageRef: ref, dockerfile: dockerfile})
		resp.Images = append(resp.Images, types.Image{
			Name: imageName,
			Tag:  ref,
		})
	}

	// Parallel runtime builds.
	bg, ctx := errgroup.WithContext(ctx)
	for _, pb := range builds {
		pb := pb

		fmt.Fprintf(params.Writer, "Building image: %s\n", pb.imageRef)

		localStart := time.Now()
		localOut := prefixer.WrapWriterWithPrefixer(params.Writer, pb.name, localStart)
		bg.Go(func() error {
			err := b.Client.BuildImage(ctx, pb.dockerfile, params.Context, pb.imageRef, excludePatterns, buildArgs, localOut)
			if err != nil {
				return err
			}

			fmt.Fprintf(params.Writer, "Built %s image in %s\n", pb.imageRef, time.Since(start).Round(time.Second))
			return nil
		})
	}
	if err := bg.Wait(); err != nil {
		return resp, err
	}

	if params.NoPush {
		return resp, nil
	}

	// Prepare pushes (skip compile).
	type pendingPush struct {
		Name string
		Tag  string // full "registry/repo:tag"
	}
	var pushes []pendingPush
	for _, buildImage := range resp.Images {
		pushes = append(pushes, pendingPush{
			Name: params.Registry,
			Tag:  image.Tag(params.Version, buildImage.Name),
		})
	}

	// Parallel pushes.
	pg, ctx := errgroup.WithContext(context.TODO())
	for _, p := range pushes {
		p := p

		fmt.Fprintf(params.Writer, "Pushing image: %s:%s\n", p.Name, p.Tag)
		out := prefixer.WrapWriterWithPrefixer(params.Writer, "push "+p.Name, start)

		pg.Go(func() error {
			localStart := time.Now()
			err := b.Client.PushImage(ctx, p.Name, p.Tag, out)
			if err != nil {
				return err
			}

			fmt.Fprintf(params.Writer, "Pushed %s:%s image in %s\n", p.Name, p.Tag, time.Since(localStart).Round(time.Second))
			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return resp, err
	}

	fmt.Fprintf(params.Writer, "Build complete in: %s\n", time.Since(start).Round(time.Second))
	return resp, nil
}

// Loads and returns a list of ignore file patterns from the specified file.
func loadIgnoreFilePatterns(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to open ignore file: %w", err)
	}
	defer f.Close()

	patterns, err := ignorefile.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing ignore file: %w", err)
	}

	return patterns, nil
}
