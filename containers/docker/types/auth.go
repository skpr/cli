package types

// Auth is a consistent object for docker authentication (package, pull) as each client does it differently.
type Auth struct {
	Username string
	Password string
	Session  string
}
