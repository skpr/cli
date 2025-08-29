package http

import (
	"fmt"
	"net/http"
	"time"
)

// Wait for a URL to become ready.
func Wait(url string, timeout time.Duration) error {
	ch := make(chan bool)
	go func() {
		for {
			_, err := http.Get(url)
			if err == nil {
				ch <- true
			}

			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("server did not reply after %v", timeout)
	}
}
