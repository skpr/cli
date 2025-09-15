package ssh

import (
	"github.com/skpr/cli/internal/client/config"
	"github.com/skpr/cli/internal/client/credentials"
	"io"
)

// Interface for the SSH client.
type Interface interface {
	Exec(ExecParams) error
	Shell(ShellParams) error
}

// ExecParams are passed to the Exec() function.
type ExecParams struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	Environment string
	Command     []string
}

// ShellParams are passed to the Shell() function.
type ShellParams struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	Environment string
}

// Client for SSH interactions.
type Client struct {
	Config      config.Config
	Credentials credentials.Credentials
}
