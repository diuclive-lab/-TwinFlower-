// Package httputil provides shared HTTP client utilities with sensible defaults.
package httputil

import (
	"net"
	"net/http"
	"time"
)

// NewClient creates an *http.Client with the given timeout and production defaults.
// Transport: 30s TLS handshake, 30s response header timeout,
// 100 max idle connections, 10 per-host idle connections.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}

// DefaultClient creates an *http.Client with a 10-second timeout.
func DefaultClient() *http.Client {
	return NewClient(10 * time.Second)
}
