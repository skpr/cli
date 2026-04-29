package logs

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
	faithcolor "github.com/fatih/color"
)

// PrettyPrint returns a colourised representation of a JSON message, falling
// back to the raw string when the message is not JSON.
func PrettyPrint(message string, indent bool) string {
	var obj map[string]interface{}

	err := json.Unmarshal([]byte(message), &obj)
	if err != nil {
		return message
	}

	formatter := colorjson.NewFormatter()
	formatter.KeyColor = faithcolor.New(faithcolor.FgWhite).Add(faithcolor.Bold)

	if indent {
		formatter.Indent = 2
	}

	raw, err := formatter.Marshal(obj)
	if err != nil {
		return message
	}

	return string(raw)
}
