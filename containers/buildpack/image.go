package buildpack

// Image build has been built.
type Image struct {
	// Name of the image.
	Name string `json:"name"`
	// Tag used to push image.
	Tag string `json:"tag"`
}

// Reference returns the image reference.
func (in *Image) Reference() string {
	// Known issue with naming of struct attributes. Keep Reference for downstream.
	return in.Tag
}
