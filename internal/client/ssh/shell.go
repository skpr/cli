package ssh

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/olekukonko/ts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Shell creates a long lived "shell" session for the user.
func (c Client) Shell(params ShellParams) error {
	awsCreds, err := c.CredentialsProvider.Retrieve(context.TODO())
	if err != nil {
		return err
	}

	var (
		user = fmt.Sprintf("%s%s%s", c.Config.Project, UsernameSeparator, params.Environment)
		pass = getPassword(awsCreds.AccessKeyID, awsCreds.SecretAccessKey, awsCreds.SessionToken)
		host = fmt.Sprintf("%s:%d", c.Cluster.SSH.Host, c.Cluster.SSH.Port)
	)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}

	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO: 1,
	}

	fd := int(os.Stdin.Fd())

	var (
		termWidth  = 80
		termHeight = 24
	)

	if terminal.IsTerminal(fd) {
		oldState, err := terminal.MakeRaw(fd)
		if err != nil {
			return err
		}

		defer terminal.Restore(fd, oldState)

		size, err := ts.GetSize()
		if err == nil {
			termWidth = size.Col()
			termHeight = size.Row()
		}
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	go monitorTerminalSizeChange(session)

	return session.Wait()
}
