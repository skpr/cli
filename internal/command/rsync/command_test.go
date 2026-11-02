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

	args := c.generateArgs("skpr", "")

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

	args := c.generateArgs("skpr", "")

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
	args := c.generateArgs("skpr", "")
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
	args := c.generateArgs("skpr", "")
	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"--dry-run",
		"source",
		"destination",
	}, args)
}

func TestGenerateArgsPrivateExclude(t *testing.T) {
	c := Command{
		Source:      "dev:/mnt/",
		Destination: "destination",
		Excludes:    []string{"js/*"},
	}
	args := c.generateArgs("skpr", "/private/")
	assert.Equal(t, []string{
		"-avz",
		"--progress",
		"-e", "skpr-rsh",
		"--exclude", "/private/",
		"--exclude", "js/*",
		"dev:/mnt/",
		"destination",
	}, args)
}

func TestPlanPrivateFiles(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		includePrivate bool
		want           privateFilesPlan
		wantErr        error
	}{
		{name: "upload", source: "./private/", want: privateFilesPlan{}},
		{name: "public download", source: "dev:/mnt/public/", want: privateFilesPlan{environment: "dev"}},
		{name: "relative download", source: "dev:files/", want: privateFilesPlan{environment: "dev"}},
		{name: "similar name", source: "dev:/mnt/private-old/", want: privateFilesPlan{environment: "dev"}},
		{name: "private without flag", source: "dev:/mnt/private/", wantErr: ErrPrivateFilesNotIncluded},
		{name: "private subdirectory without flag", source: "dev:/mnt/private/webform", wantErr: ErrPrivateFilesNotIncluded},
		{name: "private via traversal without flag", source: "dev:/mnt/public/../private", wantErr: ErrPrivateFilesNotIncluded},
		{name: "private glob without flag", source: "dev:/mnt/*", wantErr: ErrPrivateFilesNotIncluded},
		{name: "parent ending in dot without flag", source: "dev:/mnt/.", wantErr: ErrPrivateFilesNotIncluded},
		{name: "private with flag", source: "dev:/mnt/private/", includePrivate: true, want: privateFilesPlan{environment: "dev", confirm: true}},
		{name: "parent with trailing slash", source: "dev:/mnt/", want: privateFilesPlan{environment: "dev", exclude: "/private/"}},
		{name: "parent without trailing slash", source: "dev:/mnt", want: privateFilesPlan{environment: "dev", exclude: "/mnt/private/"}},
		{name: "root", source: "dev:/", want: privateFilesPlan{environment: "dev", exclude: "/mnt/private/"}},
		{name: "parent with flag", source: "dev:/mnt/", includePrivate: true, want: privateFilesPlan{environment: "dev", confirm: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Source:         tt.source,
				Destination:    "destination",
				IncludePrivate: tt.includePrivate,
			}
			got, err := c.planPrivateFiles()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
