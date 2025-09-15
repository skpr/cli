package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// EnvCluster is used as the standard way to connect to different clusters via environment variable.
	EnvCluster = "SKPR_CLUSTER"
	// EnvProject is used as the standard way to connect to different projects via environment variable.
	EnvProject = "SKPR_PROJECT"
	// EnvAPIPort is used to connect to a non standard API port.
	EnvAPIPort = "SKPR_CLUSTER_API_PORT"
	// EnvAPIInsecure is used to connect to insecure API endpoints.
	EnvAPIInsecure = "SKPR_CLUSTER_API_INSECURE"
	// EnvSSHPort is used to connect to a non standard SSH port.
	EnvSSHPort = "SKPR_CLUSTER_SSH_PORT"
)

func GetFromEnv(config *Config) error {
	var (
		cluster     = os.Getenv(EnvCluster)
		project     = os.Getenv(EnvProject)
		apiPort     = os.Getenv(EnvAPIPort)
		apiInsecure = os.Getenv(EnvAPIInsecure)
		sshPort     = os.Getenv(EnvSSHPort)
	)

	if cluster != "" {
		config.Cluster = cluster
	}

	if project != "" {
		config.Project = project
	}

	if apiPort != "" {
		val, err := strconv.Atoi(apiPort)
		if err != nil {
			return fmt.Errorf("invalid api port: %w", err)
		}

		config.API.Port = val
	}

	if apiInsecure == "true" {
		config.API.Insecure = true
	}

	if sshPort != "" {
		val, err := strconv.Atoi(sshPort)
		if err != nil {
			return fmt.Errorf("invalid ssh port: %w", err)
		}

		config.SSH.Port = val
	}

	return nil
}
