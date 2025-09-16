package pull

import "github.com/skpr/api/pb"

// Helper function to check if the image exists.
func imageExists(images []*pb.ImageStatus, id string) bool {
	for _, image := range images {
		if image.ID == id {
			return true
		}
	}

	return false
}
