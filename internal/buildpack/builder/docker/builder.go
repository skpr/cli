package buildpack

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/build"
	imagetypes "github.com/docker/docker/api/types/image"
	dockregistry "github.com/docker/docker/api/types/registry"
	dockclient "github.com/docker/docker/client"
	"github.com/moby/go-archive"
	"github.com/moby/moby/api/types/jsonstream"
	"github.com/moby/patternmatcher/ignorefile"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/buildpack/types"
	"github.com/skpr/cli/internal/buildpack/utils/finder"
	"github.com/skpr/cli/internal/buildpack/utils/image"
	"github.com/skpr/cli/internal/buildpack/utils/prefixer"
)

// Builder is the docker image builder.
type Builder struct {
	dockerClient *dockclient.Client
}

// NewBuilder creates a new Builder.
func NewBuilder() (*Builder, error) {
	dc, err := dockclient.NewClientWithOpts(dockclient.FromEnv, dockclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Builder{
		dockerClient: dc,
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

	// Build args (Docker SDK expects map[string]*string)
	buildArgs := map[string]*string{}
	withArg := func(k, v string) {
		val := v
		buildArgs[k] = &val
	}
	withArg(types.BuildArgVersion, params.Version)
	for k, v := range params.BuildArgs {
		withArg(k, v)
	}

	start := time.Now()

	// Build the compile image first; it's the base for others.
	compileRef := image.Name(params.Registry, params.Version, types.ImageNameCompile)

	fmt.Fprintf(params.Writer, "Building image: %s\n", compileRef)

	localOut := prefixer.WrapWriterWithPrefixer(params.Writer, types.ImageNameCompile, start)

	if err := b.buildImage(
		ctx,
		params.Context,
		excludePatterns,
		build.ImageBuildOptions{
			Tags:       []string{compileRef},
			Dockerfile: compileDockerfile,
			Remove:     true,
			BuildArgs:  buildArgs,
		},
		localOut,
	); err != nil {
		return resp, err
	}
	fmt.Fprintf(params.Writer, "Built %s image in %s\n", compileRef, time.Since(start).Round(time.Second))

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, types.ImageNameCompile)

	// Adds compile image identifier to the runtime images as an arg.
	withArg(types.BuildArgCompileImage, image.Name(params.Registry, params.Version, types.ImageNameCompile))

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
			err := b.buildImage(
				ctx,
				params.Context,
				excludePatterns,
				build.ImageBuildOptions{
					Tags:       []string{pb.imageRef},
					Dockerfile: pb.dockerfile,
					Remove:     true,
					BuildArgs:  buildArgs,
				},
				localOut,
			)
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
		name string
		ref  string // full "registry/repo:tag"
	}
	var pushes []pendingPush
	for _, buildImage := range resp.Images {
		pushes = append(pushes, pendingPush{
			name: buildImage.Name,
			ref:  fmt.Sprintf("%s:%s", params.Registry, image.Tag(params.Version, buildImage.Tag)),
		})
	}
	auth := dockregistry.AuthConfig{
		Username: params.Auth.Username,
		Password: params.Auth.Password,
	}
	authHdr, err := encodeRegistryAuth(auth)
	if err != nil {
		return resp, fmt.Errorf("failed to encode registry auth: %w", err)
	}

	// Parallel pushes.
	pg, ctx := errgroup.WithContext(context.TODO())
	for _, p := range pushes {
		p := p
		fmt.Fprintf(params.Writer, "Pushing image: %s\n", p.ref)
		out := prefixer.WrapWriterWithPrefixer(params.Writer, "push "+p.name, start)

		pg.Go(func() error {
			localStart := time.Now()
			rc, err := b.dockerClient.ImagePush(ctx, p.ref, imagetypes.PushOptions{
				RegistryAuth: authHdr,
			})
			if err != nil {
				return err
			}
			defer rc.Close()

			err = handleMessages(rc, out)
			if err != nil {
				return err
			}

			fmt.Fprintf(params.Writer, "Pushed %s image in %s\n", p.ref, time.Since(localStart).Round(time.Second))
			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return resp, err
	}

	fmt.Fprintf(params.Writer, "Build complete in: %s\n", time.Since(start).Round(time.Second))
	return resp, nil
}

// buildOne creates a tar build context from contextDir and streams the build output to out.
func (b *Builder) buildImage(ctx context.Context, contextDir string, contextExclude []string, opts build.ImageBuildOptions, out io.Writer) error {
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: contextExclude,
	})
	if err != nil {
		return fmt.Errorf("failed to archive build context: %w", err)
	}

	resp, err := b.dockerClient.ImageBuild(ctx, buildCtx, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = handleMessages(resp.Body, out)
	if err != nil {
		return err
	}

	return nil
}

// encodeRegistryAuth converts a registrytypes.AuthConfig into the base64-encoded JSON
// expected by the Docker Engine API for ImagePush.
func encodeRegistryAuth(cfg dockregistry.AuthConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func handleMessages(in io.ReadCloser, out io.Writer) error {
	// Stream the daemon's JSON log stream to the provided writer.
	decoder := json.NewDecoder(in)

	for {
		var msg jsonstream.Message
		if err := decoder.Decode(&msg); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// TODO: Possibly stream out progress.

		if msg.Stream != "" {
			if _, err := io.WriteString(out, msg.Stream); err != nil {
				return err
			}
		}
		if msg.Error != nil {
			return msg.Error
		}
	}

	return nil
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
