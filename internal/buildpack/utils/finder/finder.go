package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skpr/cli/internal/buildpack/builder"
)

// dockerfileSuffix is the suffix for legacy dockerfiles.
const dockerfileSuffix = ".dockerfile"

// FindDockerfiles finds the dockerfiles in the specified directory.
func FindDockerfiles(packageDir string) (buildpack.Dockerfiles, error) {
	dockerfiles := make(buildpack.Dockerfiles)
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
	if err != nil {
		return dockerfiles, err
	}
	return dockerfiles, nil
}
