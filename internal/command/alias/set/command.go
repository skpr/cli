package set

import (
	"fmt"
	"path/filepath"

	cmdconfig "github.com/skpr/cli/internal/client/config/user"
)

// Command struct.
type Command struct {
	Dir       string
	Alias     string
	Expansion string
}

// Run the command.
func (cmd *Command) Run() error {
	configPath := filepath.Join(cmd.Dir, "config.yml")
	configFile := cmdconfig.NewConfigFile(configPath)
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
