package dockerclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/build"
	imagetypes "github.com/docker/docker/api/types/image"
	dockregistry "github.com/docker/docker/api/types/registry"
	dockclient "github.com/docker/docker/client"
	"github.com/moby/go-archive"
	"github.com/moby/moby/api/types/jsonstream"
	"github.com/pkg/errors"

	"github.com/skpr/cli/internal/auth"
)

type Client struct {
	Auth   auth.Auth
	Client *dockclient.Client
}

func New(auth auth.Auth) (*Client, error) {
	client, err := dockclient.NewClientWithOpts(dockclient.FromEnv, dockclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup Docker client")
	}

	return &Client{
		Auth:   auth,
		Client: client,
	}, nil
}

func (c *Client) ImageId(ctx context.Context, name string) (string, error) {
	resp, err := c.Client.ImageInspect(ctx, name)
	if err == nil {
		return resp.ID, nil
	} else if !errdefs.IsNotFound(err) {
		return "", err
	}

	return "", nil
}

func (c *Client) PullImage(ctx context.Context, repository, tag string, writer io.Writer) error {
	auth := dockregistry.AuthConfig{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}
	authHdr, err := EncodeRegistryAuth(auth)
	if err != nil {
		return errors.Wrap(err, "failed to encode registry auth")
	}

	imageName := fmt.Sprintf("%s:%s", repository, tag)

	rc, err := c.Client.ImagePull(ctx, imageName, imagetypes.PullOptions{
		RegistryAuth: authHdr,
	})
	if err != nil {
		return err
	}

	err = handleMessages(rc, writer)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) PushImage(ctx context.Context, repository, tag string, writer io.Writer) error {
	auth := dockregistry.AuthConfig{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}
	authHdr, err := EncodeRegistryAuth(auth)
	if err != nil {
		return errors.Wrap(err, "failed to encode registry auth")
	}

	imageName := fmt.Sprintf("%s:%s", repository, tag)

	rc, err := c.Client.ImagePush(ctx, imageName, imagetypes.PushOptions{
		RegistryAuth: authHdr,
	})
	if err != nil {
		return err
	}

	err = handleMessages(rc, writer)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveImage(ctx context.Context, id string) error {
	_, err := c.Client.ImageRemove(ctx, id, imagetypes.RemoveOptions{
		PruneChildren: true,
		Force:         false,
	})
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *Client) BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error {
	buildCtx, err := archive.TarWithOptions(dockerContext, &archive.TarOptions{
		ExcludePatterns: excludePatterns,
	})
	if err != nil {
		return fmt.Errorf("failed to archive build context: %w", err)
	}

	// Build args (Docker SDK expects map[string]*string)
	args := map[string]*string{}
	withArg := func(k, v string) {
		val := v
		args[k] = &val
	}
	for k, v := range buildArgs {
		withArg(k, v)
	}

	opts := build.ImageBuildOptions{
		Tags:       []string{ref},
		Dockerfile: dockerfile,
		Remove:     true,
		BuildArgs:  args,
	}

	resp, err := c.Client.ImageBuild(ctx, buildCtx, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = handleMessages(resp.Body, writer)
	if err != nil {
		return err
	}

	return nil
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
