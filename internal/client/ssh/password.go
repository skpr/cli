package ssh

import (
	"strings"
)

// Helper function to return a combined password string.
func getPassword(username, password, session string) string {
	pass := []string{
		username,
		password,
	}

	if session != "" {
		pass = append(pass, session)
	}

	return strings.Join(pass, PasswordSeparator)
}
