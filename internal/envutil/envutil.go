// Package envutil provides typed helpers for reading environment variables.
package envutil

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String returns the env var value or the fallback if empty.
func String(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

// Int returns the env var as int or the fallback if missing/invalid.
func Int(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// Bool returns the env var as bool or the fallback if missing.
// Accepts: 1, true, yes, on → true; 0, false, no, off → false.
func Bool(key string, fallback bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

// Duration returns the env var as duration or the fallback if missing/invalid.
func Duration(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
