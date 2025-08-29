package cache

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/oauth2"
)

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
func (c Credentials) GetToken(ctx context.Context) (*oauth2.Token, error) {
	token := oauth2.Token{
		RefreshToken: c.Token.Refresh,
	}

	tokenSource := c.Config.TokenSource(ctx, &token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	return newToken, nil
}

// Cognito details for these credentials.
type Cognito struct {
	Region             string
	IdentityPoolID     string
	IdentityProviderID string
}
