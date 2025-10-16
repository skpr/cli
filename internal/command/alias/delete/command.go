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
	configFile, err := cmdconfig.NewConfigFile()
	if err != nil {
		return fmt.Errorf("could not get user config file: %w", err)
	}

	exists, err := configFile.Exists()
	if err != nil {
		return err
	}

	if !exists {
		fmt.Printf("Alias '%s' does not exist\n", cmd.Alias)
		return nil
	}

	var config cmdconfig.Config

	config, err = configFile.Read()
	if err != nil {
		return err
	}

	_, ok := config.Aliases[cmd.Alias]
	if !ok {
		fmt.Printf("Alias '%s' does not exist\n", cmd.Alias)
		return nil
	}

	delete(config.Aliases, cmd.Alias)

	err = configFile.Write(config)
	if err != nil {
		return err
	}

	fmt.Println("Alias deleted.")

	return nil
}
