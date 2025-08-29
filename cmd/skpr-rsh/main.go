package main

import (
	"os"

	wfclient "github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/ssh"
)

// The application provides the required interface for the "--rsh" rsync flag.
//
//	rsync -avz --progress --rsh 'skpr-rsh' file dev:/mnt/private/
//
// Unfortunately we cannot pass "skpr exec" because rsync passes:
//
//	skpr exec dev rsync --server -vvvvvvvlogDtprze.iLsfxC . /mnt/private/
//
// This would work if rsync passed the following:
//
//	skpr exec dev -- rsync --server -vvvvvvvlogDtprze.iLsfxC . /mnt/private/
//
// We don't have control over the format of the rsh command being passed to us
// so we have to create a shim.
//
// See the "skpr rsync" command for how this is implemented.
func main() {
	client, _, err := wfclient.NewFromFile()
	if err != nil {
		panic(err)
	}

	err = client.SSH().Exec(ssh.ExecParams{
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Environment: os.Args[1],
		Command:     os.Args[2:],
	})
	if err != nil {
		panic(err)
	}
}
