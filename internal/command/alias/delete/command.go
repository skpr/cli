package delete

import (
	"context"
	"fmt"

	cmdconfig "github.com/skpr/cli/internal/client/config/user"
)

// Command struct.
type Command struct {
	Alias string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	configFile, err := cmdconfig.NewClient()
	if err != nil {
		return fmt.Errorf("could not get user config file: %w", err)
	}

	err = configFile.RemoveAlias(cmd.Alias)
	if err != nil {
		return fmt.Errorf("failed to remove alias: %w", err)
	}

	fmt.Println("Alias deleted.")

	return nil
}
