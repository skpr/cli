package config

import (
	"fmt"
	"github.com/skpr/cli/internal/client/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	// ConfigFileName for discovery.
	ConfigFileName = "config.yml"
)

type File struct {
	Cluster string `yaml:"cluster"`
	Project string `yaml:"project"`
}

func GetFromFile(config *Config) error {
	projectDir := utils.FindSkprConfigDir(".") // @todo, Swap the dot for something better.

	if projectDir == "" {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(projectDir, ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	var file File

	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if file.Cluster != "" {
		config.Cluster = file.Cluster
	}

	if file.Project != "" {
		config.Project = file.Project
	}

	return nil
}
