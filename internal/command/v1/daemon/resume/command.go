package resume

import (
	"fmt"
	"os"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
)

// Command that resumes daemons.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	_, err = client.Daemon().Resume(ctx, &pb.DaemonResumeRequest{
		Environment: cmd.Environment,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "All Daemon tasks have been resumed.")

	return nil
}
