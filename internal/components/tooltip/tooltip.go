package tooltip

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/jwalton/gchalk"

	"github.com/skpr/cli/internal/color"
)

// Render a table to the console.
func Render(w io.Writer, msg string) error {
	var b bytes.Buffer

	t := table.New(&b)

	t.SetHeaderStyle(table.StyleBold)
	t.SetLineStyle(table.StyleBrightBlack)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetAvailableWidth(120)
	t.SetColumnMaxWidth(120)

	t.SetPadding(4)

	t.SetHeaders(gchalk.WithHex(color.HexYellow).Bold("Tooltip"))

	t.AddRow(msg)

	t.Render()

	_, err := fmt.Fprintf(w, "%s\n", b.String())
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	return nil
}
