package rsync

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/skpr/cli/internal/confirmation"
)

// ErrPrivateFilesNotIncluded is returned when a sync would download private files without the flag to allow it.
var ErrPrivateFilesNotIncluded = errors.New("this sync would download private files, which may contain personal information. " +
	"Only download private files to investigate a specific issue, with approval. " +
	"Use --include-private to download them")

// Command to rsync files.
type Command struct {
	Source         string
	Destination    string
	Excludes       []string
	ExcludeFrom    string
	DryRun         bool
	IncludePrivate bool
	Force          bool
}

// Run the command.
func (c *Command) Run(_ context.Context) error {
	executable, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to lookup binary path")
	}

	plan, err := c.planPrivateFiles()
	if err != nil {
		return err
	}

	if plan.confirm {
		msg := fmt.Sprintf("You are about to download private files from %s. Private files may contain personal information. "+
			"Only download them to investigate a specific issue, with approval. Continue? [yes/no]", plan.environment)
		if ok := confirmation.Confirm(c.Force, msg); !ok {
			return errors.New("aborted")
		}
	}

	if plan.exclude != "" {
		fmt.Fprintf(os.Stderr, "Excluding private files (%s). Use --include-private to download them.\n", PrivateFilesPath)
	}

	args := c.generateArgs(executable, plan.exclude)
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// privateFilesPlan describes how a sync handles private files.
type privateFilesPlan struct {
	// environment the files are downloaded from.
	environment string
	// confirm is true when the user must confirm a private file download.
	confirm bool
	// exclude is an rsync exclude pattern that skips private files, or empty when none is needed.
	exclude string
}

// planPrivateFiles decides how the sync handles private files. Private files are only
// downloaded with IncludePrivate set, and only after confirmation.
func (c *Command) planPrivateFiles() (privateFilesPlan, error) {
	env, remotePath, ok := splitRemote(c.Source)
	if !ok {
		return privateFilesPlan{}, nil
	}

	plan := privateFilesPlan{environment: env}

	switch classifyRemotePath(remotePath) {
	case privateScopeInside:
		if !c.IncludePrivate {
			return plan, ErrPrivateFilesNotIncluded
		}

		plan.confirm = true
	case privateScopeParent:
		if c.IncludePrivate {
			plan.confirm = true
		} else {
			plan.exclude = privateExcludePattern(remotePath)
		}
	}

	return plan, nil
}

// generateArgs generates the args for the rsync command.
func (c *Command) generateArgs(executable, privateExclude string) []string {
	args := []string{
		"-avz",
		"--progress",
		"-e", fmt.Sprintf("%s-rsh", executable),
	}
	if privateExclude != "" {
		args = append(args, "--exclude", privateExclude)
	}
	for _, ex := range c.Excludes {
		args = append(args, "--exclude", ex)
	}
	if c.ExcludeFrom != "" {
		args = append(args, "--exclude-from", c.ExcludeFrom)
	}
	if c.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, c.Source, c.Destination)
	return args
}
