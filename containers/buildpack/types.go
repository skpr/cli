package buildpack

import (
	"io"

	"github.com/skpr/cli/containers/docker/types"
)

// BuildResponse is returned by the build operation.
type BuildResponse struct {
	Images []Image `json:"images"`
}

// Params used for building the applications.
type Params struct {
	Auth       types.Auth
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
