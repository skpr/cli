package credentials

import (
	"errors"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/skpr/cli/internal/client/credentials/cache"
)

// ErrLoginRequired indicates that the developer needs to authenticate again
// with "skpr login". The cached session can no longer be renewed without them.
var ErrLoginRequired = errors.New("login required")

// Helper function to turn a token refresh failure into an error which we can
// present to a developer.
//
// The refresh grant only fails in a handful of ways and they mean very
// different things, so we split them apart here instead of leaking the raw
// oauth2 error all the way up to the console.
func classifyRefreshError(err error) error {
	if err == nil {
		return nil
	}

	// We never had a refresh token to begin with, so there is nothing to renew.
	if errors.Is(err, cache.ErrNoRefreshToken) {
		return ErrLoginRequired
	}

	var retrieve *oauth2.RetrieveError

	if errors.As(err, &retrieve) {
		switch retrieve.ErrorCode {
		// The refresh token is no longer usable eg. It has passed the Cognito
		// app client refresh token validity period, it was revoked by a global
		// sign out, or the user was disabled or reset their password.
		case "invalid_grant", "invalid_client", "unauthorized_client":
			return fmt.Errorf("%w: %w", ErrLoginRequired, err)
		}

		// A response without an error code is not an OAuth2 failure at all
		// eg. A 5xx or an error page returned by a proxy.
		return fmt.Errorf("could not reach the login provider: %w", err)
	}

	// Anything else is a transport level failure eg. DNS, timeouts.
	return fmt.Errorf("could not renew session: %w", err)
}
