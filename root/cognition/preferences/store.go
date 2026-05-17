// Package preferences learns user expression patterns so the system
// needs to clarify less over time.
package preferences

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const decayLambda = 0.05 // recency decay rate per day

// Pattern records one observed user expression pattern.
type Pattern struct {
	Pattern    string  `json:"pattern"`
	Resolved   string  `json:"resolved"`
	Confidence float64 `json:"confidence"`
	Count      int     `json:"observations"`
	CreatedAt  int64   `json:"created_at"`
	LastSeen   int64   `json:"last_seen"`
	Positive   int     `json:"positive"`
	Negative   int     `json:"negative"`
}

// effectiveConfidence applies recency decay and negative correction.
func (p *Pattern) effectiveConfidence() float64 {
	base := p.Confidence

	// Recency decay: older patterns get penalized (skip if unset)
	if p.LastSeen > 0 {
		daysSinceLastSeen := time.Since(time.Unix(p.LastSeen, 0)).Hours() / 24
		if daysSinceLastSeen > 0 {
			base *= math.Exp(-decayLambda * daysSinceLastSeen)
		}
	}

	// Negative correction: if user has corrected this pattern, reduce confidence
	if p.Negative > 0 {
		correction := float64(p.Negative) / float64(p.Positive+p.Negative+1)
		base *= (1 - correction*0.5)
	}

	// Observation boost with diminishing returns
	boost := float64(p.Count) * 0.05
	if boost > 0.30 {
		boost = 0.30
	}
	return base + boost
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
// Returns the preferred intent and decay-adjusted confidence.
func (s *Store) Lookup(input string) (string, float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lower := strings.ToLower(strings.TrimSpace(input))
	for _, p := range s.patterns {
		if matchPattern(lower, p.Pattern) {
			return p.Resolved, p.effectiveConfidence()
		}
	}
	return "", 0
}

// Observe records a user input → resolved intent mapping.
func (s *Store) Observe(input, resolved string, confidence float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	pattern := extractPattern(input)
	for i, p := range s.patterns {
		if p.Pattern == pattern && p.Resolved == resolved {
			s.patterns[i].Count++
			s.patterns[i].LastSeen = now
			s.patterns[i].Positive++
			s.patterns[i].Confidence = 0.7*confidence + 0.3*p.Confidence
			s.save()
			return
		}
	}
	// New pattern
	s.patterns = append(s.patterns, Pattern{
		Pattern:    pattern,
		Resolved:   resolved,
		Confidence: confidence,
		Count:      1,
		Positive:   1,
		CreatedAt:  now,
		LastSeen:   now,
	})
	s.save()
}

// Correct records a user correction: the pattern resolved to the wrong intent.
func (s *Store) Correct(input, wrongIntent, correctIntent string, confidence float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	pattern := extractPattern(input)

	// Penalize the wrong intent
	for i, p := range s.patterns {
		if p.Pattern == pattern && p.Resolved == wrongIntent {
			s.patterns[i].Negative++
			s.patterns[i].Confidence *= 0.7 // reduce confidence
			s.patterns[i].LastSeen = now
		}
	}

	// Add or boost the correct intent
	found := false
	for i, p := range s.patterns {
		if p.Pattern == pattern && p.Resolved == correctIntent {
			s.patterns[i].Count++
			s.patterns[i].Positive++
			s.patterns[i].LastSeen = now
			s.patterns[i].Confidence = 0.7*confidence + 0.3*p.Confidence
			found = true
			break
		}
	}
	if !found {
		s.patterns = append(s.patterns, Pattern{
			Pattern:    pattern,
			Resolved:   correctIntent,
			Confidence: confidence,
			Count:      1,
			Positive:   1,
			CreatedAt:  now,
			LastSeen:   now,
		})
	}
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
func matchPattern(input, pattern string) bool {
	parts := strings.Split(pattern, "{")
	if len(parts) == 1 {
		return input == pattern
	}
	prefix := strings.TrimSpace(parts[0])
	if prefix != "" && !strings.HasPrefix(input, prefix) {
		return false
	}
	last := strings.TrimSpace(parts[len(parts)-1])
	if strings.Contains(last, "}") {
		after := last[strings.LastIndex(last, "}")+1:]
		if after != "" && !strings.HasSuffix(input, after) {
			return false
		}
	}
	return true
}

// extractPattern generates a pattern string from user input.
func extractPattern(input string) string {
	knownCities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "南京", "武汉", "天津", "重庆"}
	result := input
	for _, city := range knownCities {
		result = strings.ReplaceAll(result, city, "{city}")
	}
	return result
}
