package mock

import (
	"context"
	"io"
	"sync"
)

// DockerClient provides a mock dockerclient client.
type DockerClient struct {
	BuildWg  sync.WaitGroup
	PushWg   sync.WaitGroup
	buildNum int
	pushNum  int
}

func (c *DockerClient) ImageId(ctx context.Context, name string) (string, error) {
	return "sha@111222333444555666", nil
}

func (c *DockerClient) PullImage(ctx context.Context, repository, tag string, writer io.Writer) error {
	return nil
}

func (c *DockerClient) PushImage(ctx context.Context, repository, tag string, writer io.Writer) error {
	c.PushWg.Done()
	c.pushNum++
	return nil
}

func (c *DockerClient) RemoveImage(ctx context.Context, id string) error {
	return nil
}

func (c *DockerClient) BuildImage(ctx context.Context, dockerfile string, dockerContext string, ref string, excludePatterns []string, buildArgs map[string]string, writer io.Writer) error {
	c.BuildWg.Done()
	c.buildNum++
	return nil
}

// BuildCount returns the build count.
func (c *DockerClient) BuildCount() int {
	c.BuildWg.Wait()
	return c.buildNum
}

// PushCount returns the push count.
func (c *DockerClient) PushCount() int {
	c.PushWg.Wait()
	return c.pushNum
}
