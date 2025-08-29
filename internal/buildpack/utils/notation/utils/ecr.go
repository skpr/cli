package utils

import (
	"fmt"

	notationregistry "github.com/notaryproject/notation-go/registry"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// GetNotationRepository creates notationregistry.Repository required to access artifacts in Amazon ECR for sign and verify operations.
// Inspired by https://github.com/aws/aws-signer-notation-plugin/blob/6e26cfd7711ad49b52c439e46bd1050c95865846/examples/utils/ecr.go
func GetNotationRepository(repository, token string) (notationregistry.Repository, error) {
	ref, err := registry.ParseReference(repository)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference: %w", err)
	}

	if err := ref.ValidateReferenceAsDigest(); err != nil {
		return nil, fmt.Errorf("invalid reference: %w", err)
	}

	authClient := &auth.Client{
		Credential: auth.StaticCredential(ref.Host(), auth.Credential{
			Username: "AWS",
			Password: token,
		}),
		ClientID: "skpr",
	}

	authClient.SetUserAgent("skpr")

	remoteRepo := &remote.Repository{
		Client:    authClient,
		Reference: ref,
	}

	err = remoteRepo.SetReferrersCapability(false)
	if err != nil {
		return nil, err
	}

	notationRepo := notationregistry.NewRepository(remoteRepo)

	return notationRepo, nil
}
