package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Dockerfiles the dockerclient build files.
type Dockerfiles map[string]string

// dockerfileSuffix is the suffix for legacy dockerfiles.
const dockerfileSuffix = ".dockerfile"

// FindDockerfiles finds the dockerfiles in the specified directory.
func FindDockerfiles(packageDir string) (Dockerfiles, error) {
	dockerfiles := make(Dockerfiles)
	err := filepath.Walk(packageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed accessing path %q: %w", path, err)
		}
		if !info.IsDir() {
			var imageName string
			if info.Name() == "Dockerfile" {
				// Use the parent directory name as the image name.
				imageName = filepath.Base(filepath.Dir(path))
			} else if strings.HasSuffix(info.Name(), dockerfileSuffix) {
				// BC: use the file name prefix as the image name.
				imageName = strings.TrimSuffix(info.Name(), dockerfileSuffix)
			} else {
				return nil
			}
			dockerfiles[imageName] = path
		}
		return nil
	})

	// Print deprecation notice.
	for key, path := range dockerfiles {
		if strings.HasSuffix(path, ".dockerfile") {
			fmt.Printf("[DEPRECATED] Dockerfile location %q is deprecated. Use \"%s/%s/Dockerfile\" instead.\n", path, filepath.Dir(path), key)
		}
	}

	if err != nil {
		return dockerfiles, err
	}
	return dockerfiles, nil
}
