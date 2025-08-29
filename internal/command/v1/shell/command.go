package shell

import (
	"os"

	"github.com/pkg/errors"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/ssh"
)

// Command to shell into an environment.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run() error {
	client, _, err := wfclient.NewFromFile()
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
