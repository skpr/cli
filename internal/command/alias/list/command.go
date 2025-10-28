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
	configFile, err := cmdconfig.NewClient()
	if err != nil {
		return fmt.Errorf("could not get user config file: %w", err)
	}

	aliases, err := configFile.ListAliases()
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(aliases) == 0 {
		fmt.Println("No aliases defined.")
		return nil
	}

	// Print the table...
	header := []string{
		"Name",
		"Command",
	}

	var rows [][]string

	for name, command := range aliases {
		rows = append(rows, []string{
			name,
			command,
		})
	}

	return table.Print(os.Stdout, header, rows)
}
