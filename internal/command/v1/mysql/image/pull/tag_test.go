package pull

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/skpr/api/pb"
)

// ImageListClient is a client for mocking the ImageListClient type.
type ImageListMockClient struct{}

// ImageList is a mock implementation of the ImageListClient interface.
func (c *ImageListMockClient) ImageList(ctx context.Context, in *pb.ImageListRequest, opts ...grpc.CallOption) (*pb.ImageListResponse, error) {
	return &pb.ImageListResponse{
		List: []*pb.ImageStatus{
			{
				ID: "skpr-rocks",
			},
		},
	}, nil
}

func TestGetTag(t *testing.T) {
	client := &ImageListMockClient{}

	// Test that a specific tag is provided.
	tag, err := getTag(context.TODO(), client, "dev", Params{
		Tag: "test-tag",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-tag", tag)

	// Test that we return an error when a tag does not exist.
	tag, err = getTag(context.TODO(), client, "default", Params{ID: "foo"})
	assert.Error(t, err)

	// Test that we return an error when a tag does not exist.
	tag, err = getTag(context.TODO(), client, "default", Params{ID: "skpr-rocks"})
	assert.NoError(t, err)
	assert.Equal(t, "skpr-rocks", tag)

	// A fallback is provided when a tag and ID is not provided.
	tag, err = getTag(context.TODO(), client, "default", Params{
		Environment: "dev",
	})
	assert.NoError(t, err)
	assert.Equal(t, "dev-default-latest", tag)
}
