package config

import (
	"os"
)

const (
	// EnvAPI is used to connect to the Skpr API server.
	EnvAPI = "SKPR_API"
	// EnvSSH is used to connect to the Skpr SSH server.
	EnvSSH = "SKPR_SSH"
	// EnvProject is used as the standard way to connect to different projects via environment variable.
	EnvProject = "SKPR_PROJECT"
)

func GetFromEnv(config *Config) error {
	var (
		api     = os.Getenv(EnvAPI)
		ssh     = os.Getenv(EnvSSH)
		project = os.Getenv(EnvProject)
	)

	if api != "" {
		config.API = URI(api)
	}

	if ssh != "" {
		config.SSH = URI(ssh)
	}

	if project != "" {
		config.Project = project
	}

	return nil
}
