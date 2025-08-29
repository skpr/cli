package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPollWithErrorBudget(t *testing.T) {
	client, err := New(
		WithBackOffBase(time.Nanosecond),
		WithBackOffIncrement(time.Nanosecond),
		WithBackOffLimit(time.Second),
	)
	assert.NoError(t, err)

	// The happy path.
	err = client.Poll(context.TODO(), func() (bool, error) {
		return true, nil
	})
	assert.NoError(t, err)

	// Fail immediately.
	err = client.Poll(context.TODO(), func() (bool, error) {
		return true, fmt.Errorf("failing")
	})
	assert.Error(t, err)

	// Fail when we hit the error limit.
	err = client.Poll(context.TODO(), func() (bool, error) {
		return false, fmt.Errorf("failing")
	})
	assert.ErrorContains(t, err, "poll has reached the limit of retries")
}
