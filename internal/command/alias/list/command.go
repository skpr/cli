package list

import (
	"fmt"
	"os"

	cmdconfig "github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/table"
)

// Command struct.
type Command struct{}

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

	if !exists {
		fmt.Println("No aliases defined.")
		return nil
	}

	config, err := configFile.Read()
	if err != nil {
		return err
	}

	if len(config.Aliases) == 0 {
		fmt.Println("No aliases defined.")
		return nil
	}

	// Print the table...
	header := []string{
		"Name",
		"Command",
	}

	var rows [][]string

	for name, command := range config.Aliases {
		rows = append(rows, []string{
			name,
			command,
		})
	}

	return table.Print(os.Stdout, header, rows)
}
