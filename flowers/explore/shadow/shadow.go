// Package shadow provides Cognitive Shadow — hypothetical execution comparison
// without runtime impact. Adapted from 66's Shadow Router, stripped of cutover logic.
package shadow

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ── Types ────────────────────────────────────────────────────────────────────

// Level defines what aspect is being shadowed.
type Level string

const (
	LevelModel      Level = "model"      // different model for the same task
	LevelThreshold  Level = "threshold"  // different clarify threshold
	LevelPreference Level = "preference" // preference layer on/off
	LevelSkill      Level = "skill"      // different skill handling
)

// Trial is one shadow comparison.
type Trial struct {
	Level       Level     `json:"level"`
	Label       string    `json:"label"`       // e.g. "threshold=0.55 vs 0.45"
	Input       string    `json:"input"`
	LiveAction  string    `json:"live_action"`  // what the production path did
	ShadowAction string   `json:"shadow_action"` // what the alternative would do
	Agreement   bool      `json:"agreement"`     // did shadow agree with live?
	LiveConf    float64   `json:"live_conf,omitempty"`
	ShadowConf  float64   `json:"shadow_conf,omitempty"`
	RanAt       time.Time `json:"ran_at"`
}

// Report aggregates shadow trials for analysis.
type Report struct {
	Level       Level    `json:"level"`
	Label       string   `json:"label"`
	TotalTrials int      `json:"total_trials"`
	Agreements  int      `json:"agreements"`
	Disagreements int    `json:"disagreements"`
	AgreeRate   float64  `json:"agree_rate"`
	Trials      []Trial  `json:"trials,omitempty"`
}

// ── Runner ───────────────────────────────────────────────────────────────────

// Runner executes shadow comparisons.
type Runner struct{}

// Compare runs the live and shadow functions, records the trial.
func (r *Runner) Compare(ctx context.Context, level Level, label, input string,
	liveFn, shadowFn func(ctx context.Context) (string, float64, error)) Trial {

	liveAction, liveConf, _ := liveFn(ctx)
	shadowAction, shadowConf, _ := shadowFn(ctx)

	return Trial{
		Level:        level,
		Label:        label,
		Input:        input,
		LiveAction:   liveAction,
		ShadowAction: shadowAction,
		Agreement:    liveAction == shadowAction,
		LiveConf:     liveConf,
		ShadowConf:   shadowConf,
		RanAt:        time.Now().UTC(),
	}
}

// ── Analysis ─────────────────────────────────────────────────────────────────

// NewReport builds a report from trials.
func NewReport(level Level, label string, trials []Trial) Report {
	agreements := 0
	for _, t := range trials {
		if t.Agreement {
			agreements++
		}
	}
	total := len(trials)
	agreeRate := 0.0
	if total > 0 {
		agreeRate = float64(agreements) / float64(total) * 100
	}
	return Report{
		Level:         level,
		Label:         label,
		TotalTrials:   total,
		Agreements:    agreements,
		Disagreements: total - agreements,
		AgreeRate:     agreeRate,
		Trials:        trials,
	}
}

// Summary returns a one-line summary of the report.
func (r Report) Summary() string {
	return fmt.Sprintf("[shadow] %s/%s: agree=%.1f%% (%d/%d)",
		r.Level, r.Label, r.AgreeRate, r.Agreements, r.TotalTrials)
}

// Markdown renders the report as readable text.
func (r Report) Markdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Cognitive Shadow: %s — %s\n\n", r.Level, r.Label))
	b.WriteString(fmt.Sprintf("- Agreement rate: **%.1f%%** (%d/%d)\n", r.AgreeRate, r.Agreements, r.TotalTrials))
	b.WriteString(fmt.Sprintf("- Disagreements: %d\n\n", r.Disagreements))

	if len(r.Trials) == 0 {
		b.WriteString("_No trials recorded._\n")
		return b.String()
	}

	// Sort: disagreements first
	sorted := make([]Trial, len(r.Trials))
	copy(sorted, r.Trials)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Agreement != sorted[j].Agreement {
			return !sorted[i].Agreement // disagreements first
		}
		return sorted[i].Input < sorted[j].Input
	})

	for _, t := range sorted {
		icon := "✅"
		if !t.Agreement {
			icon = "⚠️"
		}
		b.WriteString(fmt.Sprintf("## %s %s\n\n", icon, t.Input))
		b.WriteString(fmt.Sprintf("- **Live:** %s (conf=%.2f)\n", t.LiveAction, t.LiveConf))
		b.WriteString(fmt.Sprintf("- **Shadow:** %s (conf=%.2f)\n", t.ShadowAction, t.ShadowConf))
		if !t.Agreement {
			b.WriteString("- **Verdict:** shadow disagrees — worth investigating\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("---\n")
	b.WriteString("_Shadow trials do not affect runtime decisions. They exist to surface cognitive drift._\n")
	return b.String()
}
