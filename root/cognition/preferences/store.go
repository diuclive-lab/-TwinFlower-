// Package preferences learns user expression patterns so the system
// needs to clarify less over time.
package preferences

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// Pattern records one observed user expression pattern.
type Pattern struct {
	Pattern   string  `json:"pattern"`   // e.g. "查一下{city}"
	Resolved  string  `json:"resolved"`  // e.g. "weather"
	Confidence float64 `json:"confidence"`
	Count     int     `json:"observations"`
}

// Store is a simple JSONL-backed preference store.
type Store struct {
	mu       sync.RWMutex
	patterns []Pattern
	filePath string
}

// NewStore loads preferences from a JSONL file.
func NewStore(path string) *Store {
	s := &Store{filePath: path}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var p Pattern
		if json.Unmarshal([]byte(line), &p) == nil {
			s.patterns = append(s.patterns, p)
		}
	}
	return s
}

// Lookup checks if a user input matches a known pattern.
// Returns the preferred intent and adjusted confidence, or empty string.
func (s *Store) Lookup(input string) (string, float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lower := strings.ToLower(strings.TrimSpace(input))
	for _, p := range s.patterns {
		if matchPattern(lower, p.Pattern) {
			// More observations = higher confidence boost
			boost := float64(p.Count) * 0.05
			if boost > 0.30 {
				boost = 0.30
			}
			return p.Resolved, p.Confidence + boost
		}
	}
	return "", 0
}

// Observe records a user input → resolved intent mapping.
func (s *Store) Observe(input, resolved string, confidence float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pattern := extractPattern(input)
	for i, p := range s.patterns {
		if p.Pattern == pattern && p.Resolved == resolved {
			s.patterns[i].Count++
			// EWMA update
			s.patterns[i].Confidence = 0.7*confidence + 0.3*p.Confidence
			s.save()
			return
		}
	}
	// New pattern
	s.patterns = append(s.patterns, Pattern{
		Pattern:   pattern,
		Resolved:  resolved,
		Confidence: confidence,
		Count:     1,
	})
	s.save()
}

func (s *Store) save() {
	f, _ := os.Create(s.filePath)
	if f == nil {
		return
	}
	defer f.Close()
	for _, p := range s.patterns {
		data, _ := json.Marshal(p)
		f.Write(data)
		f.Write([]byte("\n"))
	}
}

// matchPattern checks if input matches a stored pattern.
// Currently supports simple {city} slot matching.
func matchPattern(input, pattern string) bool {
	// Simple: make pattern non-greedy, treat {slots} as wildcards
	parts := strings.Split(pattern, "{")
	if len(parts) == 1 {
		return input == pattern
	}
	// Check prefix before first slot
	prefix := strings.TrimSpace(parts[0])
	if prefix != "" && !strings.HasPrefix(input, prefix) {
		return false
	}
	// Check suffix after last slot
	last := strings.TrimSpace(parts[len(parts)-1])
	if strings.Contains(last, "}") {
		after := last[strings.LastIndex(last, "}")+1:]
		if after != "" && !strings.HasSuffix(input, after) {
			return false
		}
	}
	return true
}

// extractPattern generates a pattern string from user input by replacing
// known entity types with {slots}.
func extractPattern(input string) string {
	// Simple: replace known cities with {city}
	knownCities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "南京", "武汉", "天津", "重庆"}
	result := input
	for _, city := range knownCities {
		result = strings.ReplaceAll(result, city, "{city}")
	}
	return result
}
