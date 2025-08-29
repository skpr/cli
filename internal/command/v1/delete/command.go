package delete

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/tcnksm/go-input"

	"github.com/skpr/api/pb"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/confirmation"
)

// Command to delete and environment.
type Command struct {
	Name        string
	DryRun      bool
	SkipConfirm bool
	Force       bool
}

// Run the command.
func (cmd *Command) Run() error {
	client, ctx, err := wfclient.NewFromFile()
	if err != nil {
		return err
	}

	// Get the environment to validate it exists.
	respGet, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{Name: cmd.Name})
	if err != nil {
		// Lets provide some nicer feedback rather than just bubbling up a 404.
		return errors.Errorf("Could not find an environment named '%s'. Run 'skpr list' to view a list of deployed environments before trying again", cmd.Name)
	}

	env, err := client.Environment().Get(ctx, &pb.EnvironmentGetRequest{
		Name: cmd.Name,
	})

	if err != nil {
		return err
	}

	if env.Environment.Production {
		if okay := confirmation.Confirm(cmd.Force, "Are you sure you want to PERMANENTLY DELETE this production environment? [yes/no]"); !okay {
			return nil
		}
	}

	if cmd.DryRun {
		fmt.Println("Aborted due to --dry-run flag.")
		return nil
	}

	if cmd.SkipConfirm {
		fmt.Println("Skipping confirmation due to --yes flag.")
	} else {
		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}
		query := fmt.Sprintf("Please enter '%s' to confirm (ctrl+c to abort).", respGet.Environment.Name)
		_, err := ui.Ask(query, &input.Options{
			Required: true,
			// Validate input
			ValidateFunc: func(s string) error {
				if s != respGet.Environment.Name {
					return fmt.Errorf("input must be '%s'", respGet.Environment.Name)
				}

				return nil
			},
		})
		if err != nil {
			return errors.Wrap(err, "Error encountered when confirming delete")
		}
	}

	// Get the environment to validate it exists.
	respDel, err := client.Environment().Delete(ctx, &pb.EnvironmentDeleteRequest{Name: cmd.Name})
	if err != nil {
		return errors.Wrap(err, "Could not delete environment")
	}

	// @todo in future wait for the delete call to complete (i.e. all Finalizers to complete).

	fmt.Printf("Successfully deleted environment '%s' (status: %s)\n", respGet.Environment.Name, respDel.Status)

	return nil
}
