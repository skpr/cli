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

	projectDir, err := utils.FindSkprConfigDir()
	if err != nil {
		return fmt.Errorf("failed to find project directory: %w", err)
	}

	env, err := project.LoadFromDirectory(projectDir, cmd.Environment)
	if err != nil {
		return errors.Wrap(err, "failed to load environment")
	}

	proto, err := env.Proto(cmd.Environment, cmd.Version)
	if err != nil {
		return errors.Wrap(err, "failed to build API request")
	}

	findings, violations, warnings, err := Findings(ctx, client, proto)
	if err != nil {
		return fmt.Errorf("failed to validate environment: %w", err)
	}

	err = PrintTable(os.Stdout, findings)
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	if violations {
		// Make sure we are returning a non-zero exit code.
		// We are not using the error response because this is not an error.
		return fmt.Errorf("violations found")
	}

	if warnings && !cmd.IgnoreWarnings {
		// Make sure we are returning a non-zero exit code if not ignoring warnings.
		return fmt.Errorf("warnings found")
	}

	return nil
}

// Findings of validation checks.
func Findings(ctx context.Context, client *client.Client, proto *pb.Environment) ([]*pb.EnvironmentValidateFinding, bool, bool, error) {
	resp, err := client.Environment().Validate(ctx, &pb.EnvironmentValidateRequest{
		Environment: proto,
	})
	if err != nil {
		return nil, false, false, errors.Wrap(err, "failed to validate environment")
	}

	if len(resp.Findings) == 0 {
		return nil, false, false, nil
	}

	hasViolations, hasWarnings := false, false
	for _, finding := range resp.Findings {
		switch finding.Type {
		case pb.EnvironmentValidateFinding_Violation:
			hasViolations = true
		case pb.EnvironmentValidateFinding_Warning:
			hasWarnings = true
		}
		// Early return if both violations and warnings are found.
		if hasViolations && hasWarnings {
			return resp.Findings, true, true, nil
		}
	}
	return resp.Findings, hasViolations, hasWarnings, nil
}

// PrintTable of validation findings.
func PrintTable(w io.Writer, findings []*pb.EnvironmentValidateFinding) error {
	// Nothing to report - don't print an empty table.
	if len(findings) == 0 {
		fmt.Fprintln(w, "No validation issues found.")
	}

	header := []string{
		"Group",
		"Message",
		"Type",
	}

	var rows [][]string

	for _, finding := range findings {
		row := []string{
			finding.Group,
			finding.Message,
		}

		switch finding.Type {
		case pb.EnvironmentValidateFinding_Violation:
			row = append(row, color.New(color.FgRed).Sprintf("%s", finding.Type.String()))
		case pb.EnvironmentValidateFinding_Warning:
			row = append(row, color.New(color.FgBlue).Sprintf("%s", finding.Type.String()))
		default:
			row = append(row, finding.Type.String())
		}

		rows = append(rows, row)
	}

	err := table.Print(w, header, rows)
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	return nil
}
