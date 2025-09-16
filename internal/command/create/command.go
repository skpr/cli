package create

import (
	"context"
	"fmt"
	"github.com/skpr/cli/internal/client/project"
	"github.com/skpr/cli/internal/client/utils"
	"github.com/skpr/cli/internal/command/validate"
	"io"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
	envutils "github.com/skpr/cli/internal/environment"
)

// Command for creating an environment.
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

	projectDir := utils.FindSkprConfigDir(".") // @todo, Make this dynamic.
	if projectDir == "" {
		return fmt.Errorf("could not find project directory")
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
