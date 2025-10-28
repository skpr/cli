package user

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the persistent user configuration.
type Config struct {
	Aliases      Aliases            `yaml:"aliases,omitempty"`
	Experimental ConfigExperimental `yaml:"experimental,omitempty"`
}

// ConfigExperimental holds experimental feature flags.
type ConfigExperimental struct {
	Trace bool `yaml:"trace,omitempty"`
}

// Aliases maps alias names to commands.
type Aliases map[string]string

// ConfigFile manages the config file location.
type ConfigFile struct {
	Path string
}

// NewClient returns a ConfigFile client for interacting with ~/.skpr/config.yml.
func NewClient() (*ConfigFile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &ConfigFile{
		Path: filepath.Join(homeDir, ".skpr", "config.yml"),
	}, nil
}

// load loads config from disk, or returns an empty config if none exists.
func (c *ConfigFile) load() (Config, error) {
	var cfg Config

	data, err := os.ReadFile(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// save writes config to disk, creating directories as needed.
func (c *ConfigFile) save(cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(c.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(c.Path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SetAlias creates or updates an alias and persists the change.
func (c *ConfigFile) SetAlias(name, value string) error {
	cfg, err := c.load()
	if err != nil {
		return err
	}

	if cfg.Aliases == nil {
		cfg.Aliases = Aliases{}
	}

	cfg.Aliases[name] = value

	return c.save(cfg)
}

// RemoveAlias deletes an alias if present and persists the change.
func (c *ConfigFile) RemoveAlias(name string) error {
	cfg, err := c.load()
	if err != nil {
		return err
	}

	if cfg.Aliases != nil {
		delete(cfg.Aliases, name)
	}

	return c.save(cfg)
}

// ListAliases returns all configured aliases.
func (c *ConfigFile) ListAliases() (Aliases, error) {
	cfg, err := c.load()
	if err != nil {
		return nil, err
	}

	if cfg.Aliases == nil {
		return Aliases{}, nil
	}

	return cfg.Aliases, nil
}

// LoadFeatureFlags returns the current experimental feature flags.
func (c *ConfigFile) LoadFeatureFlags() (ConfigExperimental, error) {
	cfg, err := c.load()
	if err != nil {
		return ConfigExperimental{}, err
	}

	return cfg.Experimental, nil
}
