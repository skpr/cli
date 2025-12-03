package types

import (
	"context"

	"github.com/skpr/cli/internal/buildpack/utils/finder"
)

// Builder interface for swapping out image builders.
type Builder interface {
	Build(ctx context.Context, dockerfiles finder.Dockerfiles, params Params) (BuildResponse, error)
}
