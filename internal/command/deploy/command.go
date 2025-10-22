package deploy

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/skpr/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/project"
	"github.com/skpr/cli/internal/client/utils"
	"github.com/skpr/cli/internal/command/validate"
	envutils "github.com/skpr/cli/internal/environment"
)

// Command to delete the environment.
type Command struct {
	Environment string
	Version     string
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

	if !envutils.Contains(cmd.Environment, list.Environments) {
		return fmt.Errorf("environment not found, run 'skpr create %s %s' to provision a new environment", cmd.Environment, cmd.Version)
	}

	fmt.Println("Updating environment")

	proto, err := env.Proto(cmd.Environment, cmd.Version)
	if err != nil {
		return errors.Wrap(err, "failed to build API request")
	}

	stream, err := client.Environment().Update(ctx, &pb.EnvironmentUpdateRequest{
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
			if status.Code(err) == codes.FailedPrecondition {
				fmt.Println("Environment Validation Failed.")
				fmt.Println("Below is a list of the findings using command: skpr validate")

				violations, err := validate.PrintTable(ctx, os.Stdout, client, proto)
				if err != nil {
					return fmt.Errorf("failed to print table: %w", err)
				}

				if violations {
					// Make sure we are returning a non-zero exit code.
					// We are not using the error response because this is not an error.
					return fmt.Errorf("violations found")
				}

				return nil
			}

			return fmt.Errorf("deployment failed: %w", err)
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
