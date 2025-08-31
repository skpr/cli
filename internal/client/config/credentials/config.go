package credentials

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Clusters represents a map of clusters
type Clusters map[string]Cluster

// Cluster represents a cluster
type Cluster struct {
	AWS AWS `yaml:"aws"`
}

// AWS represents the aws config
type AWS struct {
	Profile         string `yaml:"profile,omitempty"`
	AccessKeyID     string `yaml:"access_key_id,omitempty"`
	SecretAccessKey string `yaml:"secret_access_key,omitempty"`
}

// Config represents the credentials config.
type Config struct {
	file string
}

// NewConfig create a new credentials config.
func NewConfig(file string) *Config {
	return &Config{
		file: file,
	}
}

// AddOrUpdate adds or updates a cluster in the config file.
func (c *Config) AddOrUpdate(clusterKey string, cluster Cluster) error {
	// Avoid overwriting other cluster config.
	clusters, err := c.Read()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			clusters = make(Clusters)
		} else {
			return err
		}
	}
	clusters[clusterKey] = cluster
	err = c.Write(clusters)
	return err
}

// Delete adds or updates a cluster in the config file.
func (c *Config) Delete(clusterKey string) error {
	// Avoid overwriting other cluster config.
	clusters, err := c.Read()
	if err != nil {
		return err
	}
	delete(clusters, clusterKey)
	err = c.Write(clusters)
	return err
}

// Get gets the cluster config for the specified key
func (c *Config) Get(clusterKey string) (Cluster, error) {
	clusters, err := c.Read()
	if err != nil {
		return Cluster{}, err
	}
	if _, ok := clusters[clusterKey]; !ok {
		return Cluster{}, fmt.Errorf("failed to get cluster config for %s", clusterKey)
	}
	cluster := clusters[clusterKey]
	return cluster, nil
}

// Read reads config from a file.
func (c *Config) Read() (Clusters, error) {
	data, err := os.ReadFile(c.file)
	if err != nil {
		return Clusters{}, fmt.Errorf("failed to read cluster config: %w", err)
	}

	var clusters Clusters

	err = yaml.Unmarshal(data, &clusters)
	if err != nil {
		return Clusters{}, fmt.Errorf("failed to unmarshal cluster config: %w", err)
	}
	return clusters, nil
}

// Write writes config to a file.
func (c *Config) Write(clusters Clusters) error {
	data, err := yaml.Marshal(clusters)
	if err != nil {
		return fmt.Errorf("failed to marshal clusters: %w", err)
	}

	dir := filepath.Dir(c.file)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err = os.WriteFile(c.file, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write clusters to file: %w", err)
	}

	return nil
}
