package confirmation

import (
	"fmt"

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
	if force {
		return true
	}
	if message == "" {
		message = "Are you sure? [yes/no]"
	}
	var response string
	fmt.Printf("%s: ", message)
	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err)
	}
	if slice.Contains(okayResponses, response) {
		return true
	}
	return false
}
