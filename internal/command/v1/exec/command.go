package exec

import (
	"os"

	"github.com/pkg/errors"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/ssh"
)

// Command to exec a command on an environment.
type Command struct {
	Environment string
	Command     []string
}

// Run the command.
func (cmd *Command) Run() error {
	client, _, err := wfclient.NewFromFile()
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
