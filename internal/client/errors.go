package client

import "fmt"

// ProjectInitError for errors initialising the skpr client.
type ProjectInitError struct {
	// Err is the nested error.
	Err error
}

// The Error message.
func (e ProjectInitError) Error() string {
	return fmt.Sprintf("Failed to initialise the Skpr client. Are you in a project directory? See https://docs.skpr.io/setup/project/  %s", e.Err.Error())
}

// Unwrap returns the nested error.
func (e ProjectInitError) Unwrap() error {
	return e.Err
}

// CredsError for errors initialising the skpr client.
type CredsError struct {
	// Err is the nested error.
	Err error
}

// The Error message.
func (e CredsError) Error() string {
	return fmt.Sprintf("Failed to initialise the Skpr client. Have you configured credentials? See https://docs.skpr.io/getting-started/credentials/ %s", e.Err.Error())
}

// Unwrap returns the nested error.
func (e CredsError) Unwrap() error {
	return e.Err
}
