package ecr

import (
	"context"
	"fmt"
	"strings"

	skprcredentials "github.com/skpr/cli/internal/client/credentials"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/pkg/errors"

	skprcredentials "github.com/skpr/cli/internal/client/credentials"
)

// Username to pass to the Docker registry.
// https://docs.aws.amazon.com/cli/latest/reference/ecr/get-authorization-token.html
const Username = "AWS"

// IsRegistry managed by AWS ECR.
func IsRegistry(registry string) bool {
	return strings.Contains(registry, ".ecr.")
}

// UpgradeAuth to use an AWS IAM token for authentication.
// Returns official Docker SDK registry.AuthConfig.
func UpgradeAuth(ctx context.Context, url string, creds skprcredentials.Credentials) (registrytypes.AuthConfig, error) {
	var auth registrytypes.AuthConfig

	region, err := extractRegionFromURL(url)
	if err != nil {
		return auth, errors.Wrap(err, "failed to determine registry region")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			creds.Username, creds.Password, creds.Session,
		)),
	)
	if err != nil {
		return auth, fmt.Errorf("failed to get session: %w", err)
	}

	ecrClient := ecr.NewFromConfig(cfg)

	res, err := ecrClient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return auth, err
	}

	if len(res.AuthorizationData) == 0 {
		return auth, errors.New("failed get authorization data")
	}
	ad := res.AuthorizationData[0]
	if ad.AuthorizationToken == nil {
		return auth, errors.New("failed get authorization token")
	}

	password, err := decodeAuthorizationToken(*ad.AuthorizationToken)
	if err != nil {
		return auth, errors.Wrap(err, "failed to decode authorization token")
	}

	auth.Username = Username
	auth.Password = password
	if ad.ProxyEndpoint != nil {
		// Docker accepts this with the scheme; you can strip it if your code prefers host-only.
		auth.ServerAddress = *ad.ProxyEndpoint
	}

	return auth, nil
}
