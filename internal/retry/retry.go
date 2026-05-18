// Package retry provides HTTP retry with exponential backoff.
package retry

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"
)

// DoHTTP executes fn with up to attempts retries and exponential backoff.
// Retries on: 429, 5xx, network errors, timeouts. Does not retry context errors.
func DoHTTP(ctx context.Context, attempts int, fn func() (*http.Response, error)) (*http.Response, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		resp, err := fn()
		if !shouldRetry(resp, err) {
			return resp, err
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		lastErr = err
		if resp != nil {
			lastErr = fmt.Errorf("retryable http status: %s", resp.Status)
		}
		if attempt == attempts-1 {
			break
		}
		wait := backoff(attempt)
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
	return nil, lastErr
}

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		if isContextErr(err) {
			return false
		}
		if nerr, ok := err.(net.Error); ok {
			return nerr.Timeout() || nerr.Temporary()
		}
		msg := strings.ToLower(err.Error())
		return strings.Contains(msg, "timeout") ||
			strings.Contains(msg, "connection reset") ||
			strings.Contains(msg, "broken pipe")
	}
	if resp == nil {
		return false
	}
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
}

func isContextErr(err error) bool {
	return err == context.Canceled || err == context.DeadlineExceeded
}

// backoff returns exponential backoff duration for the given attempt.
func backoff(attempt int) time.Duration {
	ms := math.Min(1000*math.Pow(2, float64(attempt)), 30000)
	return time.Duration(ms) * time.Millisecond
}
