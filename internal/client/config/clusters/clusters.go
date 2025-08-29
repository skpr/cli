package clusters

import (
	"os"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultAPIPort when not provided.
	DefaultAPIPort = 443
	// DefaultSSHPort when not provided.
	DefaultSSHPort = 22
)

// Cluster which a user connects to.
type Cluster struct {
	API API `yaml:"api"`
	SSH SSH `yaml:"ssh"`
}

// API connection details.
type API struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Insecure bool   `yaml:"insecure"`
}

// SSH connection details.
type SSH struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// LoadFromFile will return a configuration from a file.
func LoadFromFile(file, name string) (Cluster, error) {
	cluster := Cluster{
		API: API{
			Host: name,
			Port: DefaultAPIPort,
		},
		SSH: SSH{
			Host: name,
			Port: DefaultSSHPort,
		},
	}

	// If the file is not found, don't worry! Use the default cluster values.
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cluster, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return cluster, errors.Wrap(err, "failed to read config")
	}

	var clusters map[string]Cluster

	err = yaml.Unmarshal(data, &clusters)
	if err != nil {
		return cluster, errors.Wrap(err, "failed to marshal config")
	}

	// If the cluster is specified, let's use that.
	if _, ok := clusters[name]; ok {
		mergo.Merge(&cluster, clusters[name], mergo.WithOverride)
	}

	err = cluster.Validate()
	if err != nil {
		return cluster, errors.Wrap(err, "validation failed")
	}

	return cluster, nil
}

// Validate the project config file.
func (c Cluster) Validate() error {
	if c.API.Host == "" {
		return errors.New("not found: api: host")
	}

	if c.API.Port == 0 {
		return errors.New("not found: api: port")
	}

	if c.SSH.Host == "" {
		return errors.New("not found: ssh: host")
	}

	if c.SSH.Port == 0 {
		return errors.New("not found: ssh: port")
	}

	return nil
}
