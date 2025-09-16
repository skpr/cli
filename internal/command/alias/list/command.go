package list

import (
	"fmt"
	"os"
	"path/filepath"

	cmdconfig "github.com/skpr/cli/internal/client/config/user"
)

// Command struct.
type Command struct {
	Dir string
}

// Run the command.
func (cmd *Command) Run() error {
	configPath := filepath.Join(cmd.Dir, "config.yml")
	configFile := cmdconfig.NewConfigFile(configPath)
	exists, err := configFile.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("No aliases defined.")
		return nil
	}
	config, err := configFile.Read()
	if err != nil {
		return err
	}
	aliases := config.Aliases
	if len(aliases) > 0 {
		for k, v := range aliases {
			fmt.Fprintf(os.Stdout, "%-6s%s\n", k+":", v)
		}
		return nil
	}
	fmt.Println("No aliases defined.")
	return nil
}
