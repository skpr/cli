package cache

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/oauth2"
)

// ErrNoRefreshToken is returned when these credentials have no refresh token
// stored, which means they cannot be renewed.
var ErrNoRefreshToken = errors.New("refresh token not found")

// Credentials that are cached locally.
type Credentials struct {
	Config    oauth2.Config
	Token     Token
	Cognito   Cognito
	Temporary aws.Credentials
}

// TemporaryCredentialsExpired is a helper function to determine if the temporary credentials have expired.
// We wrap the AWS expired function since it does not have handled for when the credentials are empty.
func (c Credentials) TemporaryCredentialsExpired() bool {
	// We don't have any credentials at all!
	if c.Temporary.AccessKeyID == "" {
		return true
	}

	return c.Temporary.Expired()
}

// Token details for these credentials.
type Token struct {
	Refresh string
}

// GetToken returns a new token from the refresh token.
//
// Note: The identity provider does not return a refresh token as a part of the
// refresh grant response. The oauth2 package carries the existing one over for
// us, so the returned token can always be stored back as is.
func (c Credentials) GetToken(ctx context.Context) (*oauth2.Token, error) {
	// The oauth2 package returns an unexported error for this, so we check up
	// front to give callers something they can identify.
	if c.Token.Refresh == "" {
		return nil, ErrNoRefreshToken
	}

	token := oauth2.Token{
		RefreshToken: c.Token.Refresh,
	}

	tokenSource := c.Config.TokenSource(ctx, &token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// Cognito details for these credentials.
type Cognito struct {
	Region             string
	IdentityPoolID     string
	IdentityProviderID string
}
