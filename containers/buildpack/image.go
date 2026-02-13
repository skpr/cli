package buildpack

import (
	"fmt"
)

// Image build has been built.
type Image struct {
	// Name of the image.
	Name string `json:"name"`
	// Tag used to push image.
	Tag string `json:"tag"`
}

// Reference returns the image reference.
func (i *Image) Reference() string {
	return fmt.Sprintf("%d:%d", i.Name, i.Tag)
}
