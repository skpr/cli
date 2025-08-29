package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Exec a command in the remote environment.
func (c Client) Exec(params ExecParams) error {
	if len(params.Command) == 0 {
		return errors.New("Command was not provided")
	}

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

	session.Stdout = params.Stdout
	session.Stderr = params.Stderr
	session.Stdin = params.Stdin

	return session.Run(strings.Join(params.Command, " "))
}
