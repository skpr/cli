package threshold

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A malformed duration is rejected before the API is reached, so connecting
// panics rather than being reached.
func TestSetThresholdInvalid(t *testing.T) {
	for _, threshold := range []string{"", "soon", "10", "-10ms"} {
		t.Run(threshold, func(t *testing.T) {
			cmd := Command{Environment: "dev", Threshold: threshold}

			var b bytes.Buffer

			err := cmd.run(context.Background(), func(context.Context) (context.Context, api, error) {
				panic("connected despite a malformed threshold")
			}, &b)

			require.Error(t, err)
			assert.Contains(t, err.Error(), threshold)
			assert.Empty(t, b.String())
		})
	}
}
