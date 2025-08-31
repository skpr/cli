package color

import (
	"github.com/jwalton/gchalk"
)

// ApplyColorToString applies color to a status string.
func ApplyColorToString(orig string) string {
	switch orig {
	case "Deployed":
		return gchalk.WithHex(HexGreen).Sprintf("%s", orig)
	case "Completed":
		return gchalk.WithHex(HexGreen).Sprintf("%s", orig)
	case "InProgress":
		return gchalk.WithHex(HexYellow).Sprintf("%s", orig)
	case "Failed":
		return gchalk.WithHex(HexRed).Sprintf("%s", orig)
	case "Unknown":
		return gchalk.WithHex(HexRed).Sprintf("%s", orig)
	// Config types.
	case "User":
		return gchalk.WithHex(HexBlue).Sprintf("%s", orig)
	case "System":
		return gchalk.WithHex(HexGreen).Sprintf("%s", orig)
	case "Overridden":
		return gchalk.WithHex(HexYellow).Sprintf("%s", orig)
	}

	return orig
}
