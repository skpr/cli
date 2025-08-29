package color

import (
	"github.com/jwalton/gchalk"
)

// ApplyColorToString applies color to a status string.
func ApplyColorToString(orig string) string {
	switch orig {
	case "Deployed":
		return gchalk.WithHex(HexGreen).Sprintf(orig)
	case "Completed":
		return gchalk.WithHex(HexGreen).Sprintf(orig)
	case "InProgress":
		return gchalk.WithHex(HexYellow).Sprintf(orig)
	case "Failed":
		return gchalk.WithHex(HexRed).Sprintf(orig)
	case "Unknown":
		return gchalk.WithHex(HexRed).Sprintf(orig)
	// Config types.
	case "User":
		return gchalk.WithHex(HexBlue).Sprintf(orig)
	case "System":
		return gchalk.WithHex(HexGreen).Sprintf(orig)
	case "Overridden":
		return gchalk.WithHex(HexYellow).Sprintf(orig)
	}

	return orig
}
