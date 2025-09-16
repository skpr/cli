package validate

import (
	"context"
	"fmt"

	"io"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/project"
	"github.com/skpr/cli/internal/client/utils"
	"github.com/skpr/cli/internal/table"
)

// Command to validate an environments configuration.
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

	projectDir := utils.FindSkprConfigDir(".") // @todo, Make this dynamic.

	env, err := project.LoadFromDirectory(projectDir, cmd.Environment)
	if err != nil {
		return errors.Wrap(err, "failed to load environment")
	}

	proto, err := env.Proto(cmd.Environment, cmd.Version)
	if err != nil {
		return errors.Wrap(err, "failed to build API request")
	}

	violations, err := PrintTable(ctx, os.Stdout, client, proto)
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

// PrintTable of validation findings.
// Used by the create and deploy commands.
func PrintTable(ctx context.Context, w io.Writer, client *client.Client, proto *pb.Environment) (bool, error) {
	resp, err := client.Environment().Validate(ctx, &pb.EnvironmentValidateRequest{
		Environment: proto,
	})
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	header := []string{
		"Group",
		"Message",
		"Type",
	}

	var rows [][]string
	var violationCount int

	for _, finding := range resp.Findings {
		row := []string{
			finding.Group,
			finding.Message,
		}

		switch finding.Type {
		case pb.EnvironmentValidateFinding_Violation:
			violationCount++
			row = append(row, color.New(color.FgRed).Sprintf("%s", finding.Type.String()))
		case pb.EnvironmentValidateFinding_Warning:
			row = append(row, color.New(color.FgBlue).Sprintf("%s", finding.Type.String()))
		}

		rows = append(rows, row)
	}

	err = table.Print(w, header, rows)
	if err != nil {
		return false, fmt.Errorf("failed to print table: %w", err)
	}

	return violationCount > 0, nil
}
