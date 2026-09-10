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
	// DirectoryPerms for the credentials cache directory.
	DirectoryPerms = 0700
	// FilePerms for the credentials cache file. These credentials include a
	// refresh token, so they are only readable by the current user.
	FilePerms = 0600
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

	err = os.MkdirAll(directory, DirectoryPerms)
	if err != nil {
		return fmt.Errorf("failed to create credentials cache directory: %w", err)
	}

	file, err := getFile(clusterName)
	if err != nil {
		return fmt.Errorf("failed to get credentials cache file: %w", err)
	}

	return os.WriteFile(file, val, FilePerms)
}

// Delete a credentials cache file for a cluster.
func Delete(clusterName string) error {
	path, err := getFile(clusterName)
	if err != nil {
		return fmt.Errorf("failed to get credentials cache file: %w", err)
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}

	return nil
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
