package shell

import (
	"context"
	"github.com/skpr/cli/internal/client/ssh"
	"os"

	"github.com/pkg/errors"

	"github.com/skpr/cli/internal/client"
)

// Command to shell into an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	_, client, err := client.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	return client.SSH().Shell(ssh.ShellParams{
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Environment: cmd.Environment,
	})
}
