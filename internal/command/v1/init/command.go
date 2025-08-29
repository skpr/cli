package init

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/skpr/cli/internal/client/config/project"
)

// Command for initializing a project.
type Command struct {
	nonInteractive bool
	dir            string
	cluster        string
	projectName    string
}

// Run the command.
func (cmd *Command) Run() error {

	reader := bufio.NewReader(os.Stdin)

	if !cmd.nonInteractive {
		fmt.Printf("This will delete existing project files and create new ones in %s. Are you sure? [y/n]", cmd.dir)
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(confirmation)
		if confirmation != "y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	cluster := cmd.cluster
	if cluster == "" {
		fmt.Print("ClusterKey [example.skpr.io]:")
		cluster, _ = reader.ReadString('\n')
		cluster = strings.TrimSpace(cluster)
		if cluster == "" {
			return errors.New("cluster is required")
		}
	}

	projectName := cmd.projectName
	if projectName == "" {
		fmt.Print("Project name:")
		projectName, _ = reader.ReadString('\n')
		projectName = strings.TrimSpace(projectName)
		if projectName == "" {
			return errors.New("project name is required")
		}
	}

	initializer := project.NewInitializer(cmd.dir, cluster, projectName)
	err := initializer.Initialize()
	if err != nil {
		return errors.Wrap(err, "failed to initialize project")
	}
	fmt.Println("Successfully initialized project files in", cmd.dir)
	return nil
}
