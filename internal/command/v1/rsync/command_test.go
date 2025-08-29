package rsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateArgs(t *testing.T) {
	c := Command{
		Source:      "source",
		Destination: "destination",
	}

	args := c.generateArgs("skpr")

	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"source",
		"destination",
	}, args)
}

func TestGenerateArgsExcludes(t *testing.T) {
	c := Command{
		Source:      "source",
		Destination: "destination",
		Excludes: []string{
			"js/*",
			"css/*",
			"styles/*",
		},
	}

	args := c.generateArgs("skpr")

	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"--exclude", "js/*",
		"--exclude", "css/*",
		"--exclude", "styles/*",
		"source",
		"destination",
	}, args)
}

func TestGenerateArgsExcludeFrom(t *testing.T) {
	c := Command{
		Source:      "source",
		Destination: "destination",
		ExcludeFrom: "foo.txt",
	}
	args := c.generateArgs("skpr")
	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"--exclude-from", "foo.txt",
		"source",
		"destination",
	}, args)
}

func TestGenerateDryRun(t *testing.T) {
	c := Command{
		Source:      "source",
		Destination: "destination",
		DryRun:      true,
	}
	args := c.generateArgs("skpr")
	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"--dry-run",
		"source",
		"destination",
	}, args)
}
