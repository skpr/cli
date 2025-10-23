package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// Directory name for Skpr project configuration
	ProjectConfigDir = ".skpr"
	// Environment variable for declaring the project directory
	EnvProjectDirectory = "SKPR_PROJECT_DIRECTORY"
)

// FindSkprConfigDir walks up directories until it finds one containing `.skpr`.
// Returns the path to that directory, or an error if not found.
func FindSkprConfigDir() (string, error) {
	// Check for environment variable override
	envProjectDirectory := os.Getenv(EnvProjectDirectory)
	if envProjectDirectory != "" {
		exists, err := checkDirectory(envProjectDirectory)
		if err != nil {
			return "", fmt.Errorf("failed to check if directory from %s exists: %w", EnvProjectDirectory, err)
		}

		if !exists {
			return "", fmt.Errorf("directory from %s does not exist or is not a directory: %s", EnvProjectDirectory, envProjectDirectory)
		}

		return envProjectDirectory, nil
	}

	// Not environment variable was provided. Let's use the current directory.
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(curDir)
	if err != nil {
		return "", err
	}

	for {
		skprPath := filepath.Join(dir, ProjectConfigDir)
		info, err := os.Stat(skprPath)
		if err == nil && info.IsDir() {
			// Found the directory containing `.skpr`
			return filepath.Join(dir, ProjectConfigDir), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			return "", fmt.Errorf("walked to the root directory and unable to find a %s config directory", ProjectConfigDir)
		}

		dir = parent
	}
}

// Checks if a directory at the given path exists.
func checkDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
