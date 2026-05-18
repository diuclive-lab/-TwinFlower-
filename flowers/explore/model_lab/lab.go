// Package model_lab provides offline model comparison — no online dual execution.
// Adapted from 66's Router A/B, limited to research-only comparison.
package model_lab

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ── Types ────────────────────────────────────────────────────────────────────

// Experiment is a single model comparison experiment.
type Experiment struct {
	Name       string    `json:"name"`
	Prompts    []string  `json:"prompts"`
	ConfigA    Config    `json:"config_a"`
	ConfigB    Config    `json:"config_b"`
	RanAt      time.Time `json:"ran_at"`
}

// Config describes one run configuration.
type Config struct {
	Label             string  `json:"label"`                         // e.g. "qwen_dense", "gemma_dense"
	ClarifyThreshold  float64 `json:"clarify_threshold,omitempty"`
	PreferenceEnabled bool    `json:"preference_enabled"`
}

// RunResult records one prompt's output for one config.
type RunResult struct {
	Prompt   string `json:"prompt"`
	Output   string `json:"output"`
	Duration int64  `json:"duration_ms"`
	Error    string `json:"error,omitempty"`
}

// Comparison compares two configs on the same prompts.
type Comparison struct {
	Experiment string       `json:"experiment"`
	ConfigA    string       `json:"config_a"`
	ConfigB    string       `json:"config_b"`
	Results    []PairResult `json:"results"`
	WinA       int          `json:"wins_a"`
	WinB       int          `json:"wins_b"`
	Ties       int          `json:"ties"`
	RanAt      time.Time    `json:"ran_at"`
}

// PairResult is one prompt's output from both configs.
type PairResult struct {
	Prompt  string `json:"prompt"`
	OutputA string `json:"output_a"`
	OutputB string `json:"output_b"`
	Winner  string `json:"winner"` // "a", "b", or "tie"
	Note    string `json:"note,omitempty"`
}

// ── Running ──────────────────────────────────────────────────────────────────

// Handler runs a prompt against a config.
type Handler func(prompt string, cfg Config) (string, error)

// RunComparison compares two configs on all prompts.
func RunComparison(name string, prompts []string, cfgA, cfgB Config, handler Handler) Comparison {
	var results []PairResult
	winsA, winsB, ties := 0, 0, 0

	for _, prompt := range prompts {
		outA, errA := handler(prompt, cfgA)
		outB, errB := handler(prompt, cfgB)

		pr := PairResult{Prompt: prompt, OutputA: outA, OutputB: outB}

		// Simple heuristic: fewer errors = better
		switch {
		case errA != nil && errB == nil:
			pr.Winner = "b"
			winsB++
			pr.Note = fmt.Sprintf("A errored: %s", errA)
		case errB != nil && errA == nil:
			pr.Winner = "a"
			winsA++
			pr.Note = fmt.Sprintf("B errored: %s", errB)
		case errA != nil && errB != nil:
			pr.Winner = "tie"
			ties++
			pr.Note = "both errored"
		default:
			// Both succeeded — mark tie for now (could use LLM-as-judge later)
			pr.Winner = "tie"
			ties++
		}
		results = append(results, pr)
	}

	return Comparison{
		Experiment: name,
		ConfigA:    cfgA.Label,
		ConfigB:    cfgB.Label,
		Results:    results,
		WinA:       winsA,
		WinB:       winsB,
		Ties:       ties,
		RanAt:      time.Now().UTC(),
	}
}

// ── Display ──────────────────────────────────────────────────────────────────

// Markdown renders the comparison.
func (c Comparison) Markdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Model Lab: %s\n\n", c.Experiment))
	b.WriteString(fmt.Sprintf("Comparing **%s** vs **%s**\n\n", c.ConfigA, c.ConfigB))
	b.WriteString(fmt.Sprintf("| | Wins |\n|---|---|\n"))
	b.WriteString(fmt.Sprintf("| %s | %d |\n", c.ConfigA, c.WinA))
	b.WriteString(fmt.Sprintf("| %s | %d |\n", c.ConfigB, c.WinB))
	b.WriteString(fmt.Sprintf("| Ties | %d |\n\n", c.Ties))

	sort.Slice(c.Results, func(i, j int) bool {
		order := map[string]int{"b": 0, "tie": 1, "a": 2}
		return order[c.Results[i].Winner] < order[c.Results[j].Winner]
	})

	for _, r := range c.Results {
		icon := "➖"
		switch r.Winner {
		case "a":
			icon = "✅ A"
		case "b":
			icon = "✅ B"
		default:
			icon = "🔵 tie"
		}
		b.WriteString(fmt.Sprintf("## %s: %s\n\n", icon, r.Prompt))
		if r.Note != "" {
			b.WriteString(fmt.Sprintf("**Note:** %s\n\n", r.Note))
		}
	}

	b.WriteString("---\n")
	b.WriteString("_Model lab comparisons are offline research. Results do not affect production routing._\n")
	return b.String()
}
