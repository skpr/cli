package ssh

import (
	"errors"
	"fmt"
)

var (
	// ErrProjectNotFound is returned when the project is not accessible.
	ErrProjectNotFound = errors.New("project not found")
	// ErrEnvironmentNotFound is returned when the environment is not accessible.
	ErrEnvironmentNotFound = errors.New("environment not found")
)

// Helper function to handle SSH server related errors.
func handleError(project, environment string, err error) error {
	if errors.Is(err, ErrProjectNotFound) {
		fmt.Printf("Project '%s' could not be found.\n", project)
		fmt.Println("Contact your Skpr support team to help determine if this is an access or misconfiguration issue")
		return nil
	}

	if errors.Is(err, ErrEnvironmentNotFound) {
		fmt.Printf("Environment '%s' could not be found.\n", environment)
		fmt.Printf("You may have to create it. You can create it using `skpr create %s <version>`.\n", environment)
		return nil
	}

	return err
}
