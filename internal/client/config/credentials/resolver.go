package credentials

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	credentialscache "github.com/skpr/cli/internal/client/config/credentials/cache"
	providercache "github.com/skpr/cli/internal/client/config/credentials/provider/cache"
)

// Resolver is the credentials resolver.
type Resolver struct {
	config       Config
	sharedConfig string
}

// NewResolver creates a new credentials resolver.
func NewResolver(config *Config, options ...func(*Resolver)) *Resolver {
	resolver := &Resolver{
		config:       *config,
		sharedConfig: awsconfig.DefaultSharedCredentialsFilename(),
	}
	for _, option := range options {
		option(resolver)
	}
	return resolver
}

// ResolveCredentials resolves the credentials.
func (r *Resolver) ResolveCredentials(clusterKey string) (aws.CredentialsProvider, error) {
	var credsProvider aws.CredentialsProvider

	if credentialscache.Exists(clusterKey) {
		return providercache.New(clusterKey), nil
	}

	cluster, err := r.config.Get(clusterKey)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Fall back to Skpr env vars if not found.
			return NewSkprEnvCredentials(), nil
		}
		return credsProvider, err
	}

	if cluster.AWS.AccessKeyID != "" && cluster.AWS.SecretAccessKey != "" {
		return credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cluster.AWS.AccessKeyID,
				SecretAccessKey: cluster.AWS.SecretAccessKey,
				SessionToken:    "",
			},
		}, nil
	}

	if cluster.AWS.Profile != "" {
		sharedConfig, err := awsconfig.LoadSharedConfigProfile(context.TODO(), cluster.AWS.Profile, func(opts *awsconfig.LoadSharedConfigOptions) {
			opts.CredentialsFiles = []string{r.sharedConfig}
		})
		if err != nil {
			if errors.Is(err, &awsconfig.SharedConfigProfileNotExistError{}) {
				// Fall back to Skpr env vars if not found.
				return NewSkprEnvCredentials(), nil
			}
			return credsProvider, err
		}

		return aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return sharedConfig.Credentials, nil
		}), nil
	}

	return NewSkprEnvCredentials(), nil
}
