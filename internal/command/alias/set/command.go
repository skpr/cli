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
	configFile, err := cmdconfig.NewConfigFile()
	if err != nil {
		return fmt.Errorf("could not get user config file: %w", err)
	}

	exists, err := configFile.Exists()
	if err != nil {
		return err
	}

	var config cmdconfig.Config

	if !exists {
		config = cmdconfig.Config{
			Aliases: cmdconfig.Aliases{
				cmd.Alias: cmd.Expansion,
			},
		}
	} else {
		config, err = configFile.Read()
		if err != nil {
			return err
		}

		if config.Aliases == nil {
			config.Aliases = cmdconfig.Aliases{}
		}

		config.Aliases[cmd.Alias] = cmd.Expansion
	}

	err = configFile.Write(config)
	if err != nil {
		return err
	}

	fmt.Println("Alias set.")

	return nil
}
