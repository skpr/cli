package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config is the command config.
type Config struct {
	Aliases Aliases `yaml:"aliases,omitempty"`
}

// Aliases is a single command alias.
type Aliases map[string]string

// ConfigFile is the config file.
type ConfigFile struct {
	Filename string
}

// NewConfigFile creates a new config file.
func NewConfigFile(filename string) *ConfigFile {
	return &ConfigFile{
		Filename: filename,
	}
}

// Read reads config from file.
func (c *ConfigFile) Read() (Config, error) {
	var config Config
	data, err := os.ReadFile(c.Filename)
	if err != nil {
		return config, fmt.Errorf("failed to read command config: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal command config: %w", err)
	}
	return config, nil
}

// Exists checks if the config file exists.
func (c *ConfigFile) Exists() (bool, error) {
	_, err := os.Lstat(c.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Write writes config to a file.
func (c *ConfigFile) Write(config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal command config: %w", err)
	}

	dir := filepath.Dir(c.Filename)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err = ioutil.WriteFile(c.Filename, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write command config to file: %w", err)
	}

	return nil
}

// ReadAliases reads the Aliases from config.
func (c *ConfigFile) ReadAliases() (Aliases, error) {
	aliases := Aliases{}
	exists, err := c.Exists()
	if err != nil {
		return Aliases{}, err
	}
	if !exists {
		return aliases, nil
	}
	cfg, err := c.Read()
	if err != nil {
		return aliases, err
	}
	// Ensure aliases are not nil.
	if cfg.Aliases == nil {
		return aliases, nil
	}
	return cfg.Aliases, nil
}
