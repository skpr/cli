package project

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Initializer initializes projects.
type Initializer struct {
	dir    string
	config Config
}

//go:embed resources/*
var resources embed.FS

// NewInitializer creates a new initializer.
func NewInitializer(dir, cluster, projectName string) *Initializer {
	return &Initializer{
		dir: dir,
		config: Config{
			Project: projectName,
			Cluster: cluster,
		},
	}
}

// Initialize initializes a project.
func (i *Initializer) Initialize() error {
	// Create the directory if it doesn't exist.
	if _, err := os.Stat(i.dir); os.IsNotExist(err) {
		err = os.MkdirAll(i.dir, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory: %s", i.dir)
		}
	}
	// Don't include 'resources' in our paths.
	resourcesDir, err := fs.Sub(resources, "resources")
	if err != nil {
		return fmt.Errorf("failed to find resourcesDir-dir resources: %w", err)
	}
	err = fs.WalkDir(resourcesDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to load asset %s: %w", path, err)
		}
		if d.IsDir() {
			// Skip if it is a directory.
			return nil
		}
		tpl, err := template.New(filepath.Base(path)).ParseFS(resourcesDir, path)
		if err != nil {
			return fmt.Errorf("failed parsing template %s: %w", path, err)
		}
		var output bytes.Buffer
		if err = tpl.Execute(&output, i.config); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", tpl.Name(), err)
		}
		err = i.writeFile(filepath.Join(i.dir, path), output.Bytes())
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (i *Initializer) writeFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrap(err, "failed to create directory")
		}
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}
