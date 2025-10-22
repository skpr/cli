package pull

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc"
)

// ImageListClient is a client for mocking the ImageListClient type.
type ImageListClient interface {
	ImageList(ctx context.Context, in *pb.ImageListRequest, opts ...grpc.CallOption) (*pb.ImageListResponse, error)
}

// Helper function to get the tag to use for the image.
func getTag(ctx context.Context, client ImageListClient, database string, params Params) (string, error) {
	// Was a specific tag provided?
	if params.Tag != "" {
		return params.Tag, nil
	}

	if params.ID != "" {
		// Validate that the tag exists.
		images, err := client.ImageList(ctx, &pb.ImageListRequest{
			Environment: params.Environment,
		})
		if err != nil {
			return "", fmt.Errorf("failed to list images: %w", err)
		}

		if !imageExists(images.List, params.ID) {
			return "", fmt.Errorf("unable to find image with ID: %s", params.ID)
		}

		return params.ID, nil
	}

	return fmt.Sprintf("%s-%s-%s", params.Environment, database, DefaultTagSuffix), nil
}
