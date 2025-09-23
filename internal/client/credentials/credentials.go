package credentials

import (
	"context"
	"github.com/skpr/cli/internal/client/config"
)

// Credentials for connecting to the Skpr API.
type Credentials struct {
	Username string
	Password string
	Session  string
}

type ConfigGetter func(context.Context, string) (Credentials, bool, error)

func New(ctx context.Context, config config.Config) (Credentials, error) {
	funcs := []ConfigGetter{
		GetFromEnv, // Stop early if we have environment variables.
		GetFromCache,
	}

	for _, f := range funcs {
		credentials, found, err := f(ctx, config.API.Host())
		if err != nil {
			return Credentials{}, err
		}

		if found {
			return credentials, nil
		}
	}

	// If not credentials were found we will leave it up to the client to handle it.
	// That allows us to support our login command and non-authenticated commands.
	return Credentials{}, nil
}
