package buildpack

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/cli/containers/buildpack/utils/finder"
	"github.com/skpr/cli/containers/docker/mockclient"
	"github.com/skpr/cli/containers/docker/types"
)

func TestBuild(t *testing.T) {

	dockerClient := mockclient.New()

	dockerFiles := make(finder.Dockerfiles)
	dockerFiles["compile"] = ".skpr/package/compile/Dockerfile"
	dockerFiles["cli"] = ".skpr/package/cli/Dockerfile"
	dockerFiles["app"] = ".skpr/package/app/Dockerfile"
	dockerFiles["web"] = ".skpr/package/web/Dockerfile"

	var b bytes.Buffer

	params := Params{
		Writer:    &b,
		Registry:  "foo",
		Version:   "222",
		Context:   "bar",
		NoPush:    false,
		BuildArgs: map[string]string{},
		Auth:      types.Auth{},
	}

	builder, err := NewBuilder(dockerClient)
	assert.NoError(t, err)

	have, err := builder.Build(context.TODO(), dockerFiles, params)
	assert.NoError(t, err)

	want := BuildResponse{
		Images: []Image{
			{
				Name: "cli",
				Tag:  "foo:222-cli",
			},
			{
				Name: "app",
				Tag:  "foo:222-app",
			},
			{
				Name: "web",
				Tag:  "foo:222-web",
			},
		},
	}

	assert.ElementsMatch(t, want.Images, have.Images)

	// Validate we only push runtime images.
	assert.Equal(t, 4, dockerClient.BuildCount())
	assert.Equal(t, 3, dockerClient.PushCount())
}
