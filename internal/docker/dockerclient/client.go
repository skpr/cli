package dockerclient

import (
	"context"
	"fmt"
	"io"

	"github.com/containerd/errdefs"
	imagetypes "github.com/docker/docker/api/types/image"
	dockregistry "github.com/docker/docker/api/types/registry"
	dockclient "github.com/docker/docker/client"
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

	// Stream the pull progress to the UI writer.
	if _, copyErr := io.Copy(writer, rc); copyErr != nil {
		_ = rc.Close()
		return copyErr
	}
	_ = rc.Close()

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
