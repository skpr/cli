package buildpack

import (
	"fmt"
	"time"
)

// Prefixer for printing our build prefix.
type Prefixer struct {
	name  string
	start time.Time
}

// New prefixer for printing out build prefix.
func newPrefixer(name string, start time.Time) Prefixer {
	return Prefixer{
		name:  name,
		start: start,
	}
}

// PrefixFunc for dynamically calculating the prefix based on time since start.
func (p Prefixer) PrefixFunc() func() string {
	return func() string {
		return fmt.Sprintf("%s %s ", time.Since(p.start).Round(time.Second), p.name)
	}
}
