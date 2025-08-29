package list

import (
	"strings"
)

// Print a list with dots.
func Print(data []string) (string, error) {
	var dashed []string

	for _, item := range data {
		dashed = append(dashed, item)
	}

	return strings.Join(dashed, "\n"), nil
}
