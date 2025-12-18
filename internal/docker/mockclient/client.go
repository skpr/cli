package mockclient

import (
	"context"
	"io"
)

// DockerClient provides a mock client.
type DockerClient struct {
	buildNum int
	pushNum  int
}

func New() *DockerClient {
	return &DockerClient{
		buildNum: 0,
		pushNum:  0,
	}
}

func (c *DockerClient) ImageId(ctx context.Context, name string) (string, error) {
	return "sha@111222333444555666", nil
}

func (c *DockerClient) PullImage(ctx context.Context, registry, tag string, writer io.Writer) error {
	return nil
}

func (c *DockerClient) PushImage(ctx context.Context, registry, tag string, writer io.Writer) error {
	c.pushNum++
	return nil
}

func (c *DockerClient) RemoveImage(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error {
	c.buildNum++
	return nil
}

// BuildCount returns the build count.
func (c *DockerClient) BuildCount() int {
	return c.buildNum
}

// PushCount returns the push count.
func (c *DockerClient) PushCount() int {
	return c.pushNum
}
