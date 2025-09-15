package delete

import (
	"context"
	"fmt"
	"path/filepath"

	cmdconfig "github.com/skpr/cli/internal/client/config/command"
)

// Command struct.
type Command struct {
	Dir   string
	Alias string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	configPath := filepath.Join(cmd.Dir, "config.yml")
	configFile := cmdconfig.NewConfigFile(configPath)

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
