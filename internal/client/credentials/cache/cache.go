package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// Directory where the credentials will be stored.
	Directory = "skpr/credentials"
)

// Set a credentials cache file for a cluster.
func Set(clusterName string, credentials Credentials) error {
	val, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	directory, err := getDirectory()
	if err != nil {
		return fmt.Errorf("failed to get credentials cache directory: %w", err)
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create credentials cache directory: %w", err)
	}

	file, err := getFile(clusterName)
	if err != nil {
		return fmt.Errorf("failed to get credentials cache file: %w", err)
	}

	return os.WriteFile(file, val, 0644)
}

// Delete a credentials cache file for a cluster.
func Delete(clusterName string) error {
	path, err := getFile(clusterName)
	if err != nil {
		return fmt.Errorf("failed to get credentials cache file: %w", err)
	}
	return os.Remove(path)
}

// Get a credentials cache file for a cluster.
func Get(clusterName string) (Credentials, bool, error) {
	var credentials Credentials

	path, err := getFile(clusterName)
	if err != nil {
		return credentials, false, fmt.Errorf("failed to get credentials cache file: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		// We don't want to return an error if the file does not exist.
		if os.IsNotExist(err) {
			return credentials, false, nil
		}

		return credentials, false, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return credentials, false, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &credentials); err != nil {
		return credentials, false, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return credentials, true, nil
}

// Helper function to get a credentials cache file for a cluster.
func getFile(clusterName string) (string, error) {
	directory, err := getDirectory()
	if err != nil {
		return "", fmt.Errorf("failed to get credentials cache directory: %w", err)
	}

	return fmt.Sprintf("%s/%s.json", directory, clusterName), nil
}

// Helper function to get the credentials cache directory.
func getDirectory() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %w", err)
	}

	return strings.Join([]string{base, Directory}, "/"), nil
}
