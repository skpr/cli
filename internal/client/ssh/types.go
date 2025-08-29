package ssh

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/skpr/cli/internal/client/config/clusters"
	"github.com/skpr/cli/internal/client/config/project"
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
	Config              project.Config
	CredentialsProvider aws.CredentialsProvider
	Cluster             clusters.Cluster
}
