package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/skpr/cli/internal/client/utils"
)

const (
	// DefaultAPIPort when not provided.
	DefaultAPIPort = 443
	// DefaultSSHPort when not provided.
	DefaultSSHPort = 22
)

const (
	// FileName for project config discovery.
	FileName = "config.yml"
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

	data, err := os.ReadFile(filepath.Join(projectDir, FileName))
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
		config.API = URI(fmt.Sprintf("%s:%d", file.Cluster, DefaultAPIPort))
		config.SSH = URI(fmt.Sprintf("%s:%d", file.Cluster, DefaultSSHPort))
	}

	if file.Project != "" {
		config.Project = file.Project
	}

	return nil
}
