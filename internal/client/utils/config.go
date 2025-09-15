package utils

import (
	"os"
	"path/filepath"
)

const (
	ProjectConfigDir = ".skpr"
)

// FindSkprConfigDir walks up directories until it finds one containing `.skpr`.
// Returns the path to that directory, or an error if not found.
func FindSkprConfigDir(startDir string) string {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}

	for {
		skprPath := filepath.Join(dir, ProjectConfigDir)
		info, err := os.Stat(skprPath)
		if err == nil && info.IsDir() {
			// Found the directory containing `.skpr`
			return filepath.Join(dir, ProjectConfigDir)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			return ""
		}

		dir = parent
	}
}
