package confirmation

import (
	"fmt"
	"io"
	"os"

	"github.com/skpr/cli/internal/slice"
)

var (
	// okayResponses contains all the acceptable 'yes' responses
	// of which the function can determine a true/false response for.
	okayResponses = []string{"y", "Y", "yes", "Yes", "YES"}
)

// Confirm will prompt the user for a yes/no answer.
// A boolean force value is also required, which will
// force the confirmation to not show or ask for input.
// It acts as a non-interactive mode when required, but
// use with caution. This confirmation prompt is here to
// protect you.
func Confirm(force bool, message string) bool {
	return confirm(os.Stdin, os.Stdout, force, message)
}

// confirm prompts on out and reads the answer from in. An empty answer, or one that
// can't be read (e.g. no input is attached), counts as "no".
func confirm(in io.Reader, out io.Writer, force bool, message string) bool {
	if force {
		return true
	}
	if message == "" {
		message = "Are you sure? [yes/no]"
	}
	var response string
	fmt.Fprintf(out, "%s: ", message)
	if _, err := fmt.Fscanln(in, &response); err != nil {
		return false
	}
	return slice.Contains(okayResponses, response)
}
