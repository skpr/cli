package mock

import (
	"sync"

	docker "github.com/fsouza/go-dockerclient"
)

// DockerClient provides a mock docker client.
type DockerClient struct {
	BuildWg  sync.WaitGroup
	PushWg   sync.WaitGroup
	buildNum int
	pushNum  int
}

// BuildImage implements the interface.
func (c *DockerClient) BuildImage(options docker.BuildImageOptions) error {
	c.BuildWg.Done()
	c.buildNum++
	return nil
}

// PushImage implements the interface.
func (c *DockerClient) PushImage(options docker.PushImageOptions, auth docker.AuthConfiguration) error {
	c.PushWg.Done()
	c.pushNum++
	return nil
}

// InspectImage implements the interface.
func (c *DockerClient) InspectImage(name string) (*docker.Image, error) {
	return &docker.Image{
		RepoDigests: []string{
			"sha@111222333444555666",
		},
	}, nil
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
