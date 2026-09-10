package credentials

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/skpr/cli/internal/client/credentials/cache"
)

func TestClassifyRefreshError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		loginRequired bool
		message       string
	}{
		{
			name:          "no error",
			err:           nil,
			loginRequired: false,
		},
		{
			name:          "refresh token missing",
			err:           cache.ErrNoRefreshToken,
			loginRequired: true,
			message:       "login required",
		},
		{
			name:          "refresh token expired or revoked",
			err:           &oauth2.RetrieveError{ErrorCode: "invalid_grant"},
			loginRequired: true,
			message:       `login required: oauth2: "invalid_grant"`,
		},
		{
			name:          "client no longer authorized",
			err:           &oauth2.RetrieveError{ErrorCode: "unauthorized_client"},
			loginRequired: true,
			message:       `login required: oauth2: "unauthorized_client"`,
		},
		{
			name:          "identity provider is unavailable",
			err:           &oauth2.RetrieveError{Response: &http.Response{Status: "503 Service Unavailable"}, Body: []byte("unavailable")},
			loginRequired: false,
			message:       "could not reach the login provider: oauth2: cannot fetch token: 503 Service Unavailable\nResponse: unavailable",
		},
		{
			name:          "wrapped retrieve error",
			err:           fmt.Errorf("failed to get tokens: %w", &oauth2.RetrieveError{ErrorCode: "invalid_grant"}),
			loginRequired: true,
		},
		{
			name:          "transport failure",
			err:           errors.New("dial tcp: lookup example.com: no such host"),
			loginRequired: false,
			message:       "could not renew session: dial tcp: lookup example.com: no such host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := classifyRefreshError(test.err)

			if test.err == nil {
				assert.NoError(t, err)
				return
			}

			assert.Equal(t, test.loginRequired, errors.Is(err, ErrLoginRequired))

			if test.message != "" {
				assert.Equal(t, test.message, err.Error())
			}
		})
	}
}

func TestCredentialsEmpty(t *testing.T) {
	assert.True(t, Credentials{}.Empty())
	assert.True(t, Credentials{Username: "test"}.Empty())
	assert.True(t, Credentials{Password: "test"}.Empty())
	assert.False(t, Credentials{Username: "test", Password: "test"}.Empty())
}
