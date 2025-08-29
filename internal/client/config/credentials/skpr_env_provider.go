package credentials

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// EnvProviderName provides a name of Env provider
const EnvProviderName = "SkprEnvProvider"

// SkprEnvProvider is a credentials provider that uses SKPR env vars.
// * Access Key ID:     SKPR_USERNAME
// * Secret Access Key: SKPR_PASSWORD
type SkprEnvProvider struct {
	aws.CredentialsProvider
	retrieved bool
}

// NewSkprEnvCredentials returns new Skpr envvar credentials.
func NewSkprEnvCredentials() aws.CredentialsProvider {
	return &SkprEnvProvider{}
}

// Retrieve retrieves the keys from the environment.
func (e *SkprEnvProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	e.retrieved = false

	id := os.Getenv("SKPR_USERNAME")
	secret := os.Getenv("SKPR_PASSWORD")

	if id == "" {
		return aws.Credentials{Source: EnvProviderName}, errors.New("username not found")
	}

	if secret == "" {
		return aws.Credentials{Source: EnvProviderName}, errors.New("password not found")
	}

	e.retrieved = true
	return aws.Credentials{
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    "",
		Source:          EnvProviderName,
		CanExpire:       false,
	}, nil
}
