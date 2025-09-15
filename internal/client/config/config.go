package config

import (
	"errors"
	"fmt"
)

const (
	// DefaultAPIPort when not provided.
	DefaultAPIPort = 443
	// DefaultSSHPort when not provided.
	DefaultSSHPort = 22
)

// Config for connecting to the Skpr API.
type Config struct {
	Cluster string
	Project string
	API     API
	SSH     SSH
}

type API struct {
	Port     int
	Insecure bool
}

type SSH struct {
	Port int
}

type ConfigGetter func(*Config) error

func New() (Config, error) {
	config := Config{
		API: API{
			Port: DefaultAPIPort,
		},
		SSH: SSH{
			Port: DefaultSSHPort,
		},
	}

	funcs := []ConfigGetter{
		GetFromFile,
		GetFromEnv, // Environment variables should be used instead of the file if set.
	}

	for _, f := range funcs {
		err := f(&config)
		if err != nil {
			return Config{}, err
		}
	}

	var errs []error

	if config.Cluster == "" {
		errs = append(errs, fmt.Errorf("no api specified"))
	}

	if config.Project == "" {
		errs = append(errs, fmt.Errorf("no project specified"))
	}

	if config.API.Port == 0 {
		errs = append(errs, fmt.Errorf("no api port specified"))
	}

	if config.SSH.Port == 0 {
		errs = append(errs, fmt.Errorf("no ssh port specified"))
	}

	if len(errs) > 0 {
		return Config{}, errors.Join(errs...)
	}

	return config, nil
}
