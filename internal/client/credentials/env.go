package credentials

import (
	"context"
	"os"
)

const (
	// EnvUsername is used to authenticate with the Skpr cluster.
	EnvUsername = "SKPR_USERNAME"
	// EnvPassword is used to authenticate with the Skpr cluster.
	EnvPassword = "SKPR_PASSWORD"
	// EnvSession is used to authenticate with the Skpr cluster.
	EnvSession = "SKPR_SESSION"
)

func GetFromEnv(_ context.Context, _ string) (Credentials, bool, error) {
	credentials := Credentials{
		Username: os.Getenv(EnvUsername),
		Password: os.Getenv(EnvPassword),
		Session:  os.Getenv(EnvSession),
	}

	if credentials.Username == "" || credentials.Password == "" {
		// Missing required fields
		return Credentials{}, false, nil
	}

	// Both username and password are set
	return credentials, true, nil
}
