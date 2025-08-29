package image

import (
	"fmt"
)

// Name of the Docker image.
func Name(registry, version, suffix string) string {
	return fmt.Sprintf("%s:%s", registry, Tag(version, suffix))
}

// Tag assigned to a Docker image.
func Tag(version, suffix string) string {
	return fmt.Sprintf("%s-%s", version, suffix)
}
