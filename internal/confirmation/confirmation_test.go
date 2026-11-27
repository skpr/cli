package confirmation

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		force bool
		want  bool
	}{
		{name: "yes", input: "yes\n", want: true},
		{name: "short yes", input: "y\n", want: true},
		{name: "no", input: "no\n", want: false},
		{name: "empty line", input: "\n", want: false},
		{name: "no input", input: "", want: false},
		{name: "force", input: "", force: true, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, confirm(strings.NewReader(tt.input), io.Discard, tt.force, ""))
		})
	}
}
