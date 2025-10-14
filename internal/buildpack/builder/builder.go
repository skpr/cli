package buildpack

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	imagetypes "github.com/docker/docker/api/types/image"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"golang.org/x/sync/errgroup"

	"github.com/egym-playground/go-prefix-writer/prefixer"

	"github.com/skpr/cli/internal/buildpack/utils/image"
	"github.com/skpr/cli/internal/color"
)

// DockerClientInterface provides an interface that allows us to test the builder.
// This mirrors the subset of the official Docker SDK we use.
type DockerClientInterface interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error)
	ImagePush(ctx context.Context, ref string, options imagetypes.PushOptions) (io.ReadCloser, error)
	ImageInspectWithRaw(ctx context.Context, image string) (dockertypes.ImageInspect, []byte, error)
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
	Auth      registrytypes.AuthConfig
	Writer    io.Writer
	Context   string
	Registry  string
	NoPush    bool
	Version   string
	BuildArgs map[string]string
	Platform  string // e.g. "linux/amd64"
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
// Typical wiring:
//
//	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
//	b := NewBuilder(cli)
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

	// Build args (Docker SDK expects map[string]*string)
	buildArgs := map[string]*string{}
	withArg := func(k, v string) {
		val := v
		buildArgs[k] = &val
	}
	withArg(BuildArgVersion, params.Version)
	for k, v := range params.BuildArgs {
		withArg(k, v)
	}

	start := time.Now()

	// Build the compile image first; it's the base for others.
	compileRef := image.Name(params.Registry, params.Version, ImageNameCompile)
	compileOut := prefixWithTime(params.Writer, ImageNameCompile, start)

	fmt.Fprintf(params.Writer, "Building image: %s\n", compileRef)
	if err := b.buildOne(
		context.Background(),
		params.Context,
		compileDockerfile,
		dockertypes.ImageBuildOptions{
			Tags:       []string{compileRef},
			Dockerfile: compileDockerfile,
			Remove:     true,
			BuildArgs:  buildArgs,
			Platform:   params.Platform,
		},
		compileOut,
	); err != nil {
		return resp, err
	}
	fmt.Fprintf(params.Writer, "Built compile image in %s\n", time.Since(start).Round(time.Second))

	resp.Images = append(resp.Images, Image{
		Name: ImageNameCompile,
		Type: ImageTypeCompile,
		Tag:  compileRef,
	})

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, ImageNameCompile)

	// Adds compile image identifier to the runtime images as an arg.
	withArg(BuildArgCompileImage, image.Name(params.Registry, params.Version, ImageNameCompile))

	// Prepare runtime builds.
	type pendingBuild struct {
		imageRef   string
		dockerfile string
	}
	var builds []pendingBuild
	for imageName, dockerfile := range dockerfiles {
		ref := image.Name(params.Registry, params.Version, imageName)
		builds = append(builds, pendingBuild{imageRef: ref, dockerfile: dockerfile})
		resp.Images = append(resp.Images, Image{
			Name: imageName,
			Type: ImageTypeRuntime,
			Tag:  ref,
		})
	}

	// Parallel runtime builds.
	bg, ctx := errgroup.WithContext(context.Background())
	for _, pb := range builds {
		pb := pb
		out := prefixWithTime(params.Writer, pb.imageRef, start)

		fmt.Fprintf(params.Writer, "Building image: %s\n", pb.imageRef)
		bg.Go(func() error {
			localStart := time.Now()
			err := b.buildOne(
				ctx,
				params.Context,
				pb.dockerfile,
				dockertypes.ImageBuildOptions{
					Tags:       []string{pb.imageRef},
					Dockerfile: pb.dockerfile,
					Remove:     true,
					BuildArgs:  buildArgs,
					Platform:   params.Platform,
				},
				out,
			)
			if err != nil {
				return err
			}
			fmt.Fprintf(params.Writer, "Built %s image in %s\n", pb.imageRef, time.Since(localStart).Round(time.Second))
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
		ref string // full "registry/repo:tag"
	}
	var pushes []pendingPush
	for imageName := range dockerfiles {
		if imageName == ImageNameCompile {
			continue
		}
		pushes = append(pushes, pendingPush{
			ref: fmt.Sprintf("%s:%s", params.Registry, image.Tag(params.Version, imageName)),
		})
	}

	authHdr, err := encodeRegistryAuth(params.Auth)
	if err != nil {
		return resp, fmt.Errorf("failed to encode registry auth: %w", err)
	}

	// Parallel pushes.
	pg, ctx := errgroup.WithContext(context.Background())
	for _, p := range pushes {
		p := p
		fmt.Fprintf(params.Writer, "Pushing image: %s\n", p.ref)
		out := prefixWithTime(params.Writer, "push "+p.ref, start)

		pg.Go(func() error {
			localStart := time.Now()
			rc, err := b.dockerClient.ImagePush(ctx, p.ref, imagetypes.PushOptions{
				RegistryAuth: authHdr,
			})
			if err != nil {
				return err
			}
			defer rc.Close()
			if _, err := io.Copy(out, rc); err != nil {
				return err
			}
			fmt.Fprintf(params.Writer, "Pushed %s image in %s\n", p.ref, time.Since(localStart).Round(time.Second))
			return nil
		})
	}
	if err := pg.Wait(); err != nil {
		return resp, err
	}

	// Populate digests for runtime images.
	var imagesOut []Image
	for _, respImage := range resp.Images {
		if respImage.Name == ImageNameCompile {
			continue
		}
		fmt.Fprintf(params.Writer, "Fetching digest for: %s\n", respImage.Name)

		inspect, _, err := b.dockerClient.ImageInspectWithRaw(context.Background(), image.Name(params.Registry, params.Version, respImage.Name))
		if err != nil {
			return resp, fmt.Errorf("failed to inspect image %q: %w", respImage.Name, err)
		}

		digest, err := getDigest(inspect.RepoDigests)
		if err != nil {
			return resp, fmt.Errorf("failed to get digest for image %q: %w", respImage.Name, err)
		}
		respImage.Digest = digest
		imagesOut = append(imagesOut, respImage)
	}
	resp.Images = imagesOut

	fmt.Fprintf(params.Writer, "Build complete in: %s\n", time.Since(start).Round(time.Second))
	return resp, nil
}

// buildOne creates a tar build context from contextDir and streams the build output to out.
func (b *Builder) buildOne(ctx context.Context, contextDir, dockerfile string, opts dockertypes.ImageBuildOptions, out io.Writer) error {
	rc, err := tarDirectory(contextDir)
	if err != nil {
		return fmt.Errorf("failed to archive build context: %w", err)
	}
	defer rc.Close()

	resp, err := b.dockerClient.ImageBuild(ctx, rc, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Stream the daemon's JSON log stream to the provided writer.
	_, copyErr := io.Copy(out, resp.Body)
	return copyErr
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

// --- helpers below ---

// encodeRegistryAuth converts a registrytypes.AuthConfig into the base64-encoded JSON
// expected by the Docker Engine API for ImagePush.
func encodeRegistryAuth(cfg registrytypes.AuthConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// tarDirectory creates a tar archive (as ReadCloser) from the given directory.
func tarDirectory(root string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	go func() {
		tw := tar.NewWriter(pw)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if rel == "." {
				return nil
			}
			hdr, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			hdr.Name = rel

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := io.Copy(tw, f); err != nil {
					return err
				}
			}
			return nil
		})

		closeErr := tw.Close()
		if err == nil {
			err = closeErr
		}
		_ = pw.CloseWithError(err)
	}()
	return pr, nil
}

// Optional: construct a real Docker client compatible with this interface.
func NewDockerSDKClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
