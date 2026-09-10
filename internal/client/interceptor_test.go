package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	skprcredentials "github.com/skpr/cli/internal/client/credentials"
)

func TestMapAuthError(t *testing.T) {
	assert.NoError(t, mapAuthError(nil))

	// Unauthenticated is the platform telling us that we need to login again.
	err := mapAuthError(status.Error(codes.Unauthenticated, "invalid credentials"))
	assert.True(t, errors.Is(err, skprcredentials.ErrLoginRequired))

	// The platform returns this as an unknown error rather than as an
	// authentication failure.
	err = mapAuthError(status.Error(codes.Unknown, "failed to extract credentials: not found: username"))
	assert.True(t, errors.Is(err, skprcredentials.ErrLoginRequired))

	// Permission denied means we are authenticated, just not allowed.
	err = mapAuthError(status.Error(codes.PermissionDenied, "not allowed"))
	assert.False(t, errors.Is(err, skprcredentials.ErrLoginRequired))

	// All other errors are passed through untouched.
	notFound := status.Error(codes.NotFound, "environment not found")
	assert.Equal(t, notFound, mapAuthError(notFound))

	plain := errors.New("something else went wrong")
	assert.Equal(t, plain, mapAuthError(plain))
}
