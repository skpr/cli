package ssh

// Separator used for distinquishing between a namespace and exec object name.
const (
	// UsernameSeparator separates a project from the environment.
	UsernameSeparator = "~"
	// PasswordSeparator separates credentials for access.
	PasswordSeparator = ":"
)
