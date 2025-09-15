package exec

import (
	"context"
	"github.com/skpr/cli/internal/client/ssh"
	"os"

	"github.com/pkg/errors"

	"github.com/skpr/cli/internal/client"
)

// Command to exec a command on an environment.
type Command struct {
	Environment string
	Command     []string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	_, client, err := client.New(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	return client.SSH().Exec(ssh.ExecParams{
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Environment: cmd.Environment,
		Command:     cmd.Command,
	})
}
