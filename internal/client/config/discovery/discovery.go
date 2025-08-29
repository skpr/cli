package discovery

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	homedir "github.com/mitchellh/go-homedir"
)

// Discovery stores the path to skpr configuration.
type Discovery struct {
	User    User
	Project Project
}

// User configuration.
type User struct {
	Directory   string `default:"~/.skpr"`
	Credentials string `default:"credentials.yml"`
	Clusters    string `default:"clusters.yml"`
}

// Project specific configuration.
type Project struct {
	Directory string `default:".skpr"`
	Config    string `default:"config.yml"`
}

const defaultProjectDir = ".skpr"

var projectRootCmd func() (string, error)

func gitProjectRootCmd() (string, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("could not find git command in path: %w", err)
	}

	rootPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(rootPath)), nil
}

func init() {
	projectRootCmd = gitProjectRootCmd
}

// New returns a populated discovery object.
func New() (*Discovery, error) {
	d := &Discovery{}

	err := envconfig.Process("SKPR", d)

	if d.Project.Directory == "" {
		d.Project.Directory = defaultProjectDir
	}

	return d, err
}

// Credentials returns the full path to the credentials config file.
func (d *Discovery) Credentials() (string, error) {
	return homedir.Expand(path.Join(d.User.Directory, d.User.Credentials))
}

// Clusters returns the full path to the cluster config file.
func (d *Discovery) Clusters() (string, error) {
	return homedir.Expand(path.Join(d.User.Directory, d.User.Clusters))
}

// Config returns the full path to the default config.yml file.
func (d *Discovery) Config() (string, error) {
	projectDir, err := d.Project.GetDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(projectDir, d.Project.Config), nil
}

// GetDirectory returns the project directory.
func (p *Project) GetDirectory() (string, error) {
	// Check if project dir exists.
	_, err := os.Stat(p.Directory)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		// Try to find the project root.
		rootPath, err := projectRootCmd()
		if err != nil {
			return "", fmt.Errorf("could not find project root: %w", err)
		}

		return filepath.Join(rootPath, p.Directory), nil
	}

	return p.Directory, nil
}
