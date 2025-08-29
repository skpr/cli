package version

import (
	"fmt"
	"io"

	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
)

// MessageInfoUnknown is for notifying the end user information could not be found.
const MessageInfoUnknown = "UNKNOWN"

// PrintParams are passed to the Print function.
type PrintParams struct {
	ClientVersion   string
	ClientBuildDate string
	ServerVersion   string
	ServerBuildDate string
}

// Print out the version information.
func Print(w io.Writer, params PrintParams) error {
	if params.ClientVersion == "" {
		return errors.New("version not found")
	}
	if params.ClientBuildDate == "" {
		return errors.New("build date not found")
	}

	if params.ServerVersion == "" {
		params.ServerVersion = MessageInfoUnknown
	}
	if params.ServerBuildDate == "" {
		params.ServerBuildDate = MessageInfoUnknown
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("Client version:", params.ClientVersion)
	table.AddRow("Client build date:", params.ClientBuildDate)
	table.AddRow("Server version:", params.ServerVersion)
	table.AddRow("Server build date:", params.ServerBuildDate)

	_, err := fmt.Fprintln(w, table)
	return err
}
