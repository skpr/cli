package rsync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateArgs(t *testing.T) {
	tests := []struct {
		name           string
		command        Command
		privateExclude string
		want           []string
	}{
		{
			name:    "defaults",
			command: Command{Source: "source", Destination: "destination"},
			want:    []string{"-avz", "--progress", "-e", "skpr-rsh", "source", "destination"},
		},
		{
			name: "all options",
			command: Command{
				Source:      "dev:/mnt/",
				Destination: "destination",
				Excludes:    []string{"js/*", "css/*"},
				ExcludeFrom: "foo.txt",
				DryRun:      true,
			},
			privateExclude: "/private/",
			// The private exclude comes first so user patterns can't re-include private files.
			want: []string{
				"-avz", "--progress", "-e", "skpr-rsh",
				"--exclude", "/private/",
				"--exclude", "js/*",
				"--exclude", "css/*",
				"--exclude-from", "foo.txt",
				"--dry-run",
				"dev:/mnt/", "destination",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.command.generateArgs("skpr", tt.privateExclude))
		})
	}
}

func TestPlanPrivateFiles(t *testing.T) {
	plan := func(c Command) (privateFilesPlan, error) {
		c.Destination = "destination"
		return c.planPrivateFiles()
	}

	t.Run("not private", func(t *testing.T) {
		for source, env := range map[string]string{
			"./private/":            "",
			"./a:b/":                "",
			"dev:/mnt/public/":      "dev",
			"dev:files/":            "dev",
			"dev:/mnt/private-old/": "dev",
			"dev:/mnt/pub*/":        "dev",
			"dev:/var/*":            "dev",
		} {
			got, err := plan(Command{Source: source})
			assert.NoError(t, err, source)
			assert.Equal(t, privateFilesPlan{environment: env}, got, source)
		}
	})

	t.Run("requires flag", func(t *testing.T) {
		for _, source := range []string{
			"dev:/mnt/private/",
			"dev:/mnt/private/webform",
			"dev:/mnt/public/../private",
			"dev:/mnt/.",
			"dev:/mnt/*",
			"dev:/mnt/priv*/",
			"dev:/mnt/private/*.pdf",
			"dev:/m*",
			"dev:/var/*/../mnt/private",
		} {
			_, err := plan(Command{Source: source})
			assert.ErrorIs(t, err, ErrPrivateFilesNotIncluded, source)

			_, err = plan(Command{Source: source, DryRun: true})
			assert.ErrorIs(t, err, ErrPrivateFilesNotIncluded, source)

			got, err := plan(Command{Source: source, IncludePrivate: true})
			assert.NoError(t, err, source)
			assert.Equal(t, privateFilesPlan{environment: "dev", confirm: true}, got, source)

			got, err = plan(Command{Source: source, IncludePrivate: true, DryRun: true})
			assert.NoError(t, err, source)
			assert.Equal(t, privateFilesPlan{environment: "dev"}, got, source)
		}
	})

	t.Run("parent", func(t *testing.T) {
		for source, exclude := range map[string]string{
			"dev:/mnt/": "/private/",
			"dev:/mnt":  "/mnt/private/",
			"dev:/":     "/mnt/private/",
		} {
			got, err := plan(Command{Source: source})
			assert.NoError(t, err, source)
			assert.Equal(t, privateFilesPlan{environment: "dev", exclude: exclude}, got, source)

			got, err = plan(Command{Source: source, IncludePrivate: true})
			assert.NoError(t, err, source)
			assert.Equal(t, privateFilesPlan{environment: "dev", confirm: true}, got, source)
		}
	})
}

func TestRunPrivateDeclined(t *testing.T) {
	var prompted bool
	c := Command{
		Source:         "dev:/mnt/private/",
		Destination:    "destination",
		IncludePrivate: true,
		confirm: func(force bool, message string) bool {
			prompted = true
			assert.Contains(t, message, "dev")
			return false
		},
	}

	// Declining returns before rsync runs.
	assert.NoError(t, c.Run(context.Background()))
	assert.True(t, prompted)
}
