package prefixer

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/egym-playground/go-prefix-writer/prefixer"

	"github.com/skpr/cli/internal/color"
)

// Prefixer for printing our build prefix.
type Prefixer struct {
	name  string
	start time.Time
}

// NewPrefixer for printing out build prefix.
func NewPrefixer(name string, start time.Time) Prefixer {
	return Prefixer{
		name:  name,
		start: start,
	}
}

func WrapWriterWithPrefixer(w io.Writer, name string, start time.Time) io.Writer {
	return prefixer.New(w, NewPrefixer(color.Wrap(strings.ToUpper(name)), start).PrefixFunc())
}

// PrefixFunc for dynamically calculating the prefix based on time since start.
func (p Prefixer) PrefixFunc() func() string {
	return func() string {
		return fmt.Sprintf("%s %s ", time.Since(p.start).Round(time.Second), p.name)
	}
}
