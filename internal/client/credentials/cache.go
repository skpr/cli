package credentials

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/pkg/errors"
	cache2 "github.com/skpr/cli/internal/client/credentials/cache"

	skprcredentials "github.com/skpr/cli/internal/aws/credentials"
)

func GetFromCache(ctx context.Context, cluster string) (Credentials, bool, error) {
	credentials, exists, err := cache2.Get(cluster)
	if err != nil {
		return Credentials{}, false, fmt.Errorf("failed to get cached credentials: %w", err)
	}

	if !exists {
		return Credentials{}, false, nil
	}

	// If our credentials have not expired, return them.
	if !credentials.TemporaryCredentialsExpired() {
		return Credentials{
			Username: credentials.Temporary.AccessKeyID,
			Password: credentials.Temporary.SecretAccessKey,
			Session:  credentials.Temporary.SessionToken,
		}, true, nil
	}

	newToken, err := credentials.GetToken(ctx)
	if err != nil {
		return Credentials{}, false, fmt.Errorf("failed to get token: %w", err)
	}

	credentials.Token = cache2.Token{
		Refresh: newToken.RefreshToken,
	}

	// Extract the ID Token from OAuth2 token.
	idToken, ok := newToken.Extra("id_token").(string)
	if !ok {
		return Credentials{}, false, errors.Wrap(err, "Missing id_token")
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(credentials.Cognito.Region), awsconfig.WithCredentialsProvider(aws.AnonymousCredentials{}))
	if err != nil {
		return Credentials{}, false, fmt.Errorf("failed to load AWS config: %w", err)
	}

	newTemporary, err := skprcredentials.GetTempCredentials(ctx, cognitoidentity.NewFromConfig(cfg), skprcredentials.GetTempCredentialsParams{
		Token:            idToken,
		IdentityPool:     credentials.Cognito.IdentityPoolID,
		IdentityProvider: credentials.Cognito.IdentityProviderID,
	})
	if err != nil {
		return Credentials{}, false, fmt.Errorf("failed to get temporary credentials: %w", err)
	}

	// Set our name so that we can identify the source of the credentials.
	newTemporary.Source = "SkprCacheProvider" // @todo, Make this a const again.

	credentials.Temporary = newTemporary

	// Save back the refreshed token and credentials.
	err = cache2.Set(cluster, credentials)
	if err != nil {
		return Credentials{}, false, fmt.Errorf("failed to store refreshed credentials: %w", err)
	}

	return Credentials{
		Username: credentials.Temporary.AccessKeyID,
		Password: credentials.Temporary.SecretAccessKey,
		Session:  credentials.Temporary.SessionToken,
	}, true, nil
}
