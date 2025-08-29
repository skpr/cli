package rsync

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Command to rsync files.
type Command struct {
	Source      string
	Destination string
	Excludes    []string
	ExcludeFrom string
	DryRun      bool
}

// Run the command.
func (c *Command) Run() error {
	executable, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to lookup binary path")
	}

	args := c.generateArgs(executable)
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// generateArgs generates the args for the rsync command.
func (c *Command) generateArgs(executable string) []string {
	args := []string{
		"-avz",
		"--progress",
		"-e", fmt.Sprintf("%s-rsh", executable),
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
