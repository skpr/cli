package errorhandler

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/fang"

	"github.com/skpr/cli/internal/client/credentials"
	"github.com/skpr/cli/internal/components/tooltip"
)

const (
	// EnvDebug can be set by a developer to see the underlying error.
	EnvDebug = "SKPR_DEBUG"
	// Message displayed when a developer needs to authenticate.
	messageLoginRequired = "You are not authenticated with Skpr, or your session has expired."
	// Tooltip displayed with the command a developer should run next.
	tooltipLoginRequired = "Authenticate with this command:\n\n$ skpr login\n"
)

// Handler renders an error and, where we know it, the command which the
// developer should run next.
func Handler(w io.Writer, styles fang.Styles, err error) {
	if errors.Is(err, credentials.ErrLoginRequired) {
		renderLoginRequired(w, styles, err)
		return
	}

	fang.DefaultErrorHandler(w, styles, err)
}

// Helper function to render the authentication error.
func renderLoginRequired(w io.Writer, styles fang.Styles, err error) {
	_, _ = fmt.Fprintln(w, styles.ErrorHeader.String())
	_, _ = fmt.Fprintln(w, styles.ErrorText.Render(messageLoginRequired))
	_, _ = fmt.Fprintln(w)

	if renderErr := tooltip.Render(w, tooltipLoginRequired); renderErr != nil {
		_, _ = fmt.Fprintln(w, styles.ErrorText.Render("Authenticate with: skpr login"))
		_, _ = fmt.Fprintln(w)
	}

	// We have replaced the underlying error with our own message, so we keep a
	// way for a developer to see what actually happened.
	if os.Getenv(EnvDebug) != "" {
		_, _ = fmt.Fprintln(w, styles.ErrorText.Render(fmt.Sprintf("Debug: %v", err)))
		_, _ = fmt.Fprintln(w)
	}
}
