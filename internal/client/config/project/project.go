package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	// FileDefaults is used for loading environment defaults.
	FileDefaults = "defaults.yml"

	// EnvConfigProject is the environment variable for the project name.
	EnvConfigProject = "SKPR_CONFIG_PROJECT"
	// EnvConfigCluster is the environment variable for the cluster name.
	EnvConfigCluster = "SKPR_CONFIG_CLUSTER"
)

// LoadFromDirectory project configuration.
func LoadFromDirectory(path, name string) (Environment, error) {
	var environment Environment

	defaults := filepath.Join(path, FileDefaults)

	// Load default environment configuration.
	if _, err := os.Stat(defaults); os.IsNotExist(err) {
		return environment, err
	}

	data, err := os.ReadFile(defaults)
	if err != nil {
		return environment, err
	}

	err = yaml.Unmarshal(data, &environment)
	if err != nil {
		return environment, err
	}

	// Load the environment specific configuration.
	data, err = os.ReadFile(filepath.Join(path, fmt.Sprintf("%s.yml", name)))
	if err != nil && !os.IsNotExist(err) {
		return environment, err
	}

	err = yaml.Unmarshal(data, &environment)
	if err != nil {
		return environment, err
	}

	return environment, nil
}

// LoadConfig from a file and fallback to environment variables.
func LoadConfig(filePath string) (Config, error) {
	var config Config

	// If the configuration file exists. Load it from that file.
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return config, fmt.Errorf("failed to get project config: %w", err)
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return config, fmt.Errorf("failed to unmarshal project config: %w", err)
		}
	}

	if os.Getenv(EnvConfigCluster) != "" {
		config.Cluster = os.Getenv(EnvConfigCluster)
	}

	if os.Getenv(EnvConfigProject) != "" {
		config.Project = os.Getenv(EnvConfigProject)
	}

	// Validate configuration.
	if config.Project == "" {
		return config, fmt.Errorf("project configuration not found")
	}

	if config.Cluster == "" {
		return config, fmt.Errorf("cluster configuration not found")
	}

	return config, nil
}
