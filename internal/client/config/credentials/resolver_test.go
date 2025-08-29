package credentials

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreds(t *testing.T) {
	config := NewConfig("testdata/.skpr/credentials.yml")
	resolver := NewResolver(config)
	credentialsProvider, err := resolver.ResolveCredentials("foo.cluster.skpr.io")
	assert.NoError(t, err)

	creds, err := credentialsProvider.Retrieve(context.TODO())
	assert.NoError(t, err)

	assert.Equal(t, "xxxx", creds.AccessKeyID)
	assert.Equal(t, "yyyy", creds.SecretAccessKey)

}

func TestProfile(t *testing.T) {
	config := NewConfig("testdata/.skpr/credentials.yml")

	profileConfig := func(res *Resolver) {
		res.sharedConfig = "testdata/credentials"
	}
	resolver := NewResolver(config, profileConfig)
	credentialsProvider, err := resolver.ResolveCredentials("bar.cluster.skpr.io")
	assert.NoError(t, err)

	creds, err := credentialsProvider.Retrieve(context.TODO())
	assert.NoError(t, err)

	assert.Equal(t, "fizzy", creds.AccessKeyID)
	assert.Equal(t, "popper", creds.SecretAccessKey)
}

func TestSkprEnvVars(t *testing.T) {

	err := os.Setenv("SKPR_USERNAME", "wiz")
	assert.NoError(t, err)
	err = os.Setenv("SKPR_PASSWORD", "bang")
	assert.NoError(t, err)

	config := NewConfig("NOT_FOUND")
	resolver := NewResolver(config)
	credentialsProvider, err := resolver.ResolveCredentials("bar.cluster.skpr.io")
	assert.NoError(t, err)

	creds, err := credentialsProvider.Retrieve(context.TODO())
	assert.NoError(t, err)

	assert.Equal(t, "wiz", creds.AccessKeyID)
	assert.Equal(t, "bang", creds.SecretAccessKey)

}
