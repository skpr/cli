package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Client for retrying function calls.
type Client struct {
	// Logger for debugging.
	Logger *slog.Logger
	// How many errors to tolerate before failing.
	ErrorLimit int32
	Timeout    time.Duration
	// How long to initial backoff for.
	BackOffBase time.Duration
	// How much we increment our backoff.
	BackOffIncrement time.Duration
	// Limit for how long we wait during a backoff event.
	BackOffLimit time.Duration
}

// Option for configuring the client.
type Option func(c *Client) error

// WithLogger sets the logger on the client.
func WithLogger(val *slog.Logger) Option {
	return func(c *Client) error {
		c.Logger = val
		return nil
	}
}

// WithErrorLimit sets the error limit on the client.
func WithErrorLimit(val int32) Option {
	return func(c *Client) error {
		c.ErrorLimit = val
		return nil
	}
}

// WithTimeout sets the timeout on the client.
func WithTimeout(val time.Duration) Option {
	return func(c *Client) error {
		c.Timeout = val
		return nil
	}
}

// WithBackOffBase sets the backoff base on the client.
func WithBackOffBase(val time.Duration) Option {
	return func(c *Client) error {
		c.BackOffBase = val
		return nil
	}
}

// WithBackOffIncrement sets the backoff increment on the client.
func WithBackOffIncrement(val time.Duration) Option {
	return func(c *Client) error {
		c.BackOffIncrement = val
		return nil
	}
}

// WithBackOffLimit sets the backoff limit on the client.
func WithBackOffLimit(val time.Duration) Option {
	return func(c *Client) error {
		c.BackOffLimit = val
		return nil
	}
}

// New client for retrying function calls.
func New(options ...Option) (*Client, error) {
	client := &Client{
		ErrorLimit:       15,
		Timeout:          30 * time.Minute,
		BackOffBase:      10 * time.Second,
		BackOffIncrement: 5 * time.Second,
		BackOffLimit:     5 * time.Minute,
	}

	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, fmt.Errorf("failed to set option: %w", err)
		}
	}

	return client, nil
}

// PollFunc that is provided to PollWithErrorBudget.
type PollFunc func() (bool, error)

// Poll and wait for function to complete.
func (c *Client) Poll(ctx context.Context, pollFunc PollFunc) error {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	// Used for checking how many times the command has failed.
	var errCount int32

	// Used for tracking time between executions.
	backoff := c.BackOffBase

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context deadline exceeded")
		default:
			exit, err := pollFunc()

			// The function has instructed this poll to exit with an error.
			if exit && err != nil {
				return err
			}

			// The function has instructed this poll to complete.
			if exit {
				return nil
			}

			// Did we return an error?
			if err != nil {
				errCount++

				if c.Logger != nil {
					c.Logger.Error(err.Error())
				}
			}

			if errCount >= c.ErrorLimit {
				return fmt.Errorf("poll has reached the limit of retries: %d", c.ErrorLimit)
			}

			backoff = getNextBackOff(backoff, c.BackOffIncrement, c.BackOffLimit)

			time.Sleep(backoff)
		}
	}
}

// Returns the next backoff duration.
func getNextBackOff(current, increment, limit time.Duration) time.Duration {
	current += increment

	if current >= limit {
		return limit
	}

	return current
}
