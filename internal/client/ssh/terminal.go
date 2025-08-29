//go:build !windows
// +build !windows

package ssh

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/olekukonko/ts"
	"golang.org/x/crypto/ssh"
)

// Helper function to update the height and width on window resize.
func monitorTerminalSizeChange(session *ssh.Session) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// resize the tty if any signals received
	for range sigs {
		size, err := ts.GetSize()
		if err == nil {
			session.WindowChange(size.Row(), size.Col())
		}
	}
}
