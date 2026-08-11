package create

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/project"
	"github.com/skpr/cli/internal/client/utils"
	"github.com/skpr/cli/internal/command/validate"
	envutils "github.com/skpr/cli/internal/environment"
)

// Command for creating an environment.
type Command struct {
	Environment    string
	Version        string
	IgnoreWarnings bool
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, client, err := client.New(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Loading environment configuration")

	projectDir, err := utils.FindSkprConfigDir()
	if err != nil {
		return fmt.Errorf("failed to find project directory: %w", err)
	}

	env, err := project.LoadFromDirectory(projectDir, cmd.Environment)
	if err != nil {
		return errors.Wrap(err, "failed to load environment")
	}

	list, err := client.Environment().List(ctx, &pb.EnvironmentListRequest{})
	if err != nil {
		return errors.Wrap(err, "failed to list environments")
	}

	if envutils.Contains(cmd.Environment, list.Environments) {
		return fmt.Errorf("environment already exists, run 'skpr deploy %s %s' to update the existing environment", cmd.Environment, cmd.Version)
	}

	fmt.Println("Creating environment")

	proto, err := env.Proto(cmd.Environment, cmd.Version)
	if err != nil {
		return errors.Wrap(err, "failed to build API request")
	}

	findings, violations, warnings, err := validate.Findings(ctx, client, proto)
	if err != nil {
		return fmt.Errorf("failed to validate environment: %w", err)
	}

	if violations || (warnings && !cmd.IgnoreWarnings) {
		err = validate.PrintTable(os.Stdout, findings)
		if err != nil {
			return fmt.Errorf("failed to print table: %w", err)
		}
		return fmt.Errorf("validation issues found")
	}

	stream, err := client.Environment().Create(ctx, &pb.EnvironmentCreateRequest{
		Environment: proto,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("environment creation failed: %w", err)
		}

		fmt.Println(resp.Message)
	}

	if _, exists := os.LookupEnv("SKPR_AWESOME_LOGS"); exists {
		fmt.Println("Now you're off to the races!")
	} else {
		fmt.Println("Complete")
	}

	return nil
}
