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
