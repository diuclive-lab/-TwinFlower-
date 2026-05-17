// Package calibration records model behavior per intent for adaptive tuning.
package calibration

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Record is one observation of model behavior for a specific intent.
type Record struct {
	Timestamp  string  `json:"timestamp"`
	Model      string  `json:"model"`
	Intent     string  `json:"intent"`
	Tool       string  `json:"tool"`
	Confidence float64 `json:"confidence"`
	Clarified  bool    `json:"clarified"`
	Success    bool    `json:"success"`
	LatencyMs  int64   `json:"latency_ms"`
}

var mu sync.Mutex

// Log appends a calibration record to runtime_stats.jsonl.
// Creates the file if it doesn't exist.
func Log(r Record) {
	r.Timestamp = time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(r)
	mu.Lock()
	defer mu.Unlock()
	f, err := os.OpenFile("runtime/calibration.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(data)
	f.Write([]byte("\n"))
}
