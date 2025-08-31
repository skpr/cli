package list

import (
	"strings"
)

// Print a list with dots.
func Print(data []string) (string, error) {
	var dashed []string
	dashed = append(dashed, data...)
	return strings.Join(dashed, "\n"), nil
}
