package set

import (
	"fmt"

	cmdconfig "github.com/skpr/cli/internal/client/config/user"
)

// Command struct.
type Command struct {
	Alias     string
	Expansion string
}

// Run the command.
func (cmd *Command) Run() error {
	configFile, err := cmdconfig.NewClient()
	if err != nil {
		return fmt.Errorf("could not get user config file: %w", err)
	}

	err = configFile.SetAlias(cmd.Alias, cmd.Expansion)
	if err != nil {
		return fmt.Errorf("failed to set alias: %w", err)
	}

	fmt.Println("Alias set.")

	return nil
}
