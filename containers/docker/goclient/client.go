package goclient

import (
	"context"
	"io"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"

	"github.com/skpr/cli/containers/docker/types"
)

type Client struct {
	Auth   types.Auth
	Client *dockerclient.Client
}

func New(auth types.Auth) (*Client, error) {
	client, err := dockerclient.NewClientFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup Docker client")
	}

	return &Client{
		Auth:   auth,
		Client: client,
	}, nil
}

func (c *Client) ImageId(ctx context.Context, name string) (string, error) {
	resp, err := c.Client.InspectImage(name)
	if err != nil && !errors.Is(err, dockerclient.ErrNoSuchImage) {
		return "", err
	}

	if resp == nil {
		return "", nil
	}

	return resp.ID, nil
}

func (c *Client) PullImage(ctx context.Context, registry, tag string, writer io.Writer) error {
	opts := dockerclient.PullImageOptions{
		OutputStream: writer,
		Repository:   registry,
		Tag:          tag,
		Context:      ctx,
	}

	clientAuth := dockerclient.AuthConfiguration{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}

	return c.Client.PullImage(opts, clientAuth)
}

func (c *Client) PushImage(ctx context.Context, registry, tag string, writer io.Writer) error {
	opts := dockerclient.PushImageOptions{
		Context:      ctx,
		OutputStream: writer,
		Name:         registry,
		Tag:          tag,
	}

	clientAuth := dockerclient.AuthConfiguration{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}

	return c.Client.PushImage(opts, clientAuth)
}

func (c *Client) RemoveImage(ctx context.Context, id string) error {
	err := c.Client.RemoveImage(id)
	return err
}

func (c *Client) BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error {
	args := []dockerclient.BuildArg{}
	for k, v := range buildArgs {
		args = append(args, dockerclient.BuildArg{
			Name:  k,
			Value: v,
		})
	}

	build := dockerclient.BuildImageOptions{
		Context:      ctx,
		Name:         ref,
		Dockerfile:   dockerfile,
		ContextDir:   dockerContext,
		OutputStream: writer,
		BuildArgs:    args,
	}
	return c.Client.BuildImage(build)
}
