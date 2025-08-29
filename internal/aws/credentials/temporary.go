package credentials

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

// GetTempCredentialsParams is provided to the GetTempCredentials function.
type GetTempCredentialsParams struct {
	Token            string
	IdentityPool     string
	IdentityProvider string
}

// Validate input provided to GetTempCredentials function.
func (p GetTempCredentialsParams) Validate() error {
	var errs []error

	if p.Token == "" {
		errs = append(errs, fmt.Errorf("token not provided"))
	}

	if p.IdentityPool == "" {
		errs = append(errs, fmt.Errorf("identity pool not provided"))
	}

	if p.IdentityProvider == "" {
		errs = append(errs, fmt.Errorf("identity provider not provided"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// CognitoClient defines the functions required when we are interacting with Cognito.
type CognitoClient interface {
	GetId(ctx context.Context, params *cognitoidentity.GetIdInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.GetIdOutput, error)
	GetCredentialsForIdentity(ctx context.Context, params *cognitoidentity.GetCredentialsForIdentityInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.GetCredentialsForIdentityOutput, error)
}

// GetTempCredentials gets the temporary STS AWS credentials for the oauth tokens, and saves them.
func GetTempCredentials(ctx context.Context, identity CognitoClient, params GetTempCredentialsParams) (aws.Credentials, error) {
	var credentials aws.Credentials

	if err := params.Validate(); err != nil {
		return credentials, fmt.Errorf("validation failed: %w", err)
	}

	logins := map[string]string{
		params.IdentityProvider: params.Token,
	}

	id, err := identity.GetId(ctx, &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(params.IdentityPool),
		Logins:         logins,
	})
	if err != nil {
		return credentials, fmt.Errorf("failed to get cognito id: %w", err)
	}

	c, err := identity.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: id.IdentityId,
		Logins:     logins,
	})
	if err != nil {
		return credentials, fmt.Errorf("failed to get credentials for identity: %w", err)
	}

	credentials.AccessKeyID = *c.Credentials.AccessKeyId
	credentials.SecretAccessKey = *c.Credentials.SecretKey
	credentials.SessionToken = *c.Credentials.SessionToken
	credentials.Expires = *c.Credentials.Expiration
	credentials.CanExpire = true

	return credentials, nil
}
