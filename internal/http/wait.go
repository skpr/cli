package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Wait for a URL to become ready.
func Wait(ctx context.Context, url string, timeout time.Duration) error {
	var (
		ch   = make(chan bool)
		done = make(chan struct{})
	)

	// Stop polling when we return, otherwise this goroutine is left running for
	// the lifetime of the command.
	defer close(done)

	go func() {
		for {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()

				select {
				case ch <- true:
				case <-done:
				}

				return
			}

			select {
			case <-time.After(1 * time.Second):
			case <-done:
				return
			}
		}
	}()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(timeout):
		return fmt.Errorf("server did not reply after %v", timeout)
	}
}
