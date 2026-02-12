package image

import (
	"fmt"
)

// Name of the Docker image.
func Name(registry, tag string) string {
	return fmt.Sprintf("%s:%s", registry, tag)
}

// Tag assigned to a Docker image.
func Tag(version, suffix string) string {
	return fmt.Sprintf("%s-%s", version, suffix)
}
