package table

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/jwalton/gchalk"

	"github.com/skpr/cli/internal/color"
)

// Print a table to the console.
func Print(w io.Writer, headers []string, rows [][]string) error {
	var b bytes.Buffer

	t := table.New(&b)

	t.SetHeaderStyle(table.StyleBold)
	t.SetLineStyle(table.StyleBrightBlack)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetAvailableWidth(80)
	t.SetColumnMaxWidth(80)

	var formattedHeaders []string

	for _, h := range headers {
		formattedHeaders = append(formattedHeaders, gchalk.WithHex(color.HexOrange).Bold(h))
	}

	t.SetHeaders(formattedHeaders...)

	for _, row := range rows {
		t.AddRow(row...)
	}

	t.Render()

	_, err := fmt.Fprintf(w, "\n%s\n", b.String())
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	return nil
}
