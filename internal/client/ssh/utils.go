package ssh

import (
	"fmt"
	"io"
	"strings"

	"github.com/gliderlabs/ssh"
)

// Separator used for distinquishing between a namespace and exec object name.
const (
	// UsernameSeparator separates a project from the environment.
	UsernameSeparator = "~"
	// PasswordSeparator separates credentials for access.
	PasswordSeparator = ":"
)

// Helper function to return errors back to the user.
func printError(s ssh.Session, err error) {
	io.WriteString(s, err.Error())
	s.Exit(1)
}

// Helper function to extract namespace and name from a user.
func extractFromUser(user string) (string, string, error) {
	sl := strings.Split(user, UsernameSeparator)

	if len(sl) == 2 {
		// Translate this
		return sl[0], sl[1], nil
	}

	return "", "", fmt.Errorf("failed to marshal string: %s", user)
}
