package buildpack

import (
	"io"

	"github.com/skpr/cli/internal/auth"
)

// BuildResponse is returned by the build operation.
type BuildResponse struct {
	Images []Image `json:"images"`
}

// Image build has been built.
type Image struct {
	// Name of the image.
	Name string `json:"name"`
	// Tag used to push image.
	Tag string `json:"tag"`
}

// Params used for building the applications.
type Params struct {
	Auth       auth.Auth
	Writer     io.Writer
	Context    string
	IgnoreFile string
	Registry   string
	NoPush     bool
	Version    string
	BuildArgs  map[string]string
}

const (
	// ImageNameCompile is used for compiling the application.
	ImageNameCompile = "compile"

	// BuildArgCompileImage is used for referencing the compile image.
	BuildArgCompileImage = "COMPILE_IMAGE"
	// BuildArgVersion is used for providing the version identifier of the application.
	BuildArgVersion = "SKPR_VERSION"
)
