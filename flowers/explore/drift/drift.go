// Package drift detects preference and behavior shifts over time.
// Adapted from 66's drift detection, repositioned as cognitive observation,
// not a release gate.
package drift

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// ── Types ────────────────────────────────────────────────────────────────────

// Kind is the type of drift being observed.
type Kind string

const (
	KindPreference Kind = "preference" // user expression → tool mapping changed
	KindThreshold  Kind = "threshold"  // clarify rate changed
	KindBehavior   Kind = "behavior"   // how the system responds changed
)

// Severity indicates how significant the drift is.
type Severity string

const (
	SeverityNone    Severity = "none"
	SeverityLow     Severity = "low"
	SeverityMedium  Severity = "medium"
	SeverityHigh    Severity = "high"
)

// Observation is one data point for drift analysis.
type Observation struct {
	Pattern   string    `json:"pattern"`
	OldIntent string    `json:"old_intent"`
	NewIntent string    `json:"new_intent,omitempty"`
	OldConf   float64   `json:"old_confidence"`
	NewConf   float64   `json:"new_confidence"`
	Count     int       `json:"count"`
	WindowA   string    `json:"window_a"` // e.g. "last_7d"
	WindowB   string    `json:"window_b"` // e.g. "previous_7d"
	Severity  Severity  `json:"severity"`
	DetectedAt time.Time `json:"detected_at"`
}

// Report aggregates drift observations.
type Report struct {
	GeneratedAt  time.Time     `json:"generated_at"`
	Observations []Observation `json:"observations"`
	AlertCount   int           `json:"alert_count"`   // severity >= medium
	ChangeCount  int           `json:"change_count"`  // intent changed (not just confidence)
}

// ── Detection ────────────────────────────────────────────────────────────────

// Detector analyzes calibration data for drift.
type Detector struct {
	// Threshold for considering a change significant (0.0-1.0).
	ConfidenceGap float64
	// Minimum observations before considering a pattern stable.
	MinObservations int
}

// NewDetector creates a detector with sensible defaults.
func NewDetector() *Detector {
	return &Detector{
		ConfidenceGap:   0.15,
		MinObservations: 3,
	}
}

// ComparePrefs checks if a preference pattern has drifted.
// oldConf and newConf are the effective confidence from the preference store.
func (d *Detector) ComparePrefs(pattern, oldIntent, newIntent string, oldConf, newConf float64, count int) Observation {
	sev := SeverityNone
	gap := math.Abs(newConf - oldConf)

	switch {
	case gap >= d.ConfidenceGap*2 && oldIntent != newIntent:
		sev = SeverityHigh
	case gap >= d.ConfidenceGap && oldIntent == newIntent:
		sev = SeverityMedium
	case gap >= d.ConfidenceGap:
		sev = SeverityMedium
	case gap >= d.ConfidenceGap*0.5:
		sev = SeverityLow
	}

	return Observation{
		Pattern:    pattern,
		OldIntent:  oldIntent,
		NewIntent:  newIntent,
		OldConf:    oldConf,
		NewConf:    newConf,
		Count:      count,
		Severity:   sev,
		DetectedAt: time.Now().UTC(),
	}
}

// ── Analysis ─────────────────────────────────────────────────────────────────

// NewReport builds a report from observations.
func NewReport(observations []Observation) Report {
	alerts := 0
	changes := 0
	for _, o := range observations {
		if o.Severity == SeverityMedium || o.Severity == SeverityHigh {
			alerts++
		}
		if o.OldIntent != "" && o.NewIntent != "" && o.OldIntent != o.NewIntent {
			changes++
		}
	}
	return Report{
		GeneratedAt:  time.Now().UTC(),
		Observations: observations,
		AlertCount:   alerts,
		ChangeCount:  changes,
	}
}

// Markdown renders the drift report.
func (r Report) Markdown() string {
	var b strings.Builder
	b.WriteString("# Preference Drift Report\n\n")

	if len(r.Observations) == 0 {
		b.WriteString("_No drift detected. All patterns stable._\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("- Alerts (severity ≥ medium): **%d**\n", r.AlertCount))
	b.WriteString(fmt.Sprintf("- Intent changes: **%d**\n\n", r.ChangeCount))

	// Sort: high severity first
	sorted := make([]Observation, len(r.Observations))
	copy(sorted, r.Observations)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Severity > sorted[j].Severity })

	for _, o := range sorted {
		if o.Severity == SeverityNone {
			continue
		}
		badge := ""
		switch o.Severity {
		case SeverityHigh:
			badge = "🔴 HIGH"
		case SeverityMedium:
			badge = "🟡 MEDIUM"
		case SeverityLow:
			badge = "🟢 LOW"
		}

		b.WriteString(fmt.Sprintf("## %s: %s\n\n", badge, o.Pattern))
		b.WriteString(fmt.Sprintf("- **Old behavior:** %s (conf=%.2f)\n", o.OldIntent, o.OldConf))
		b.WriteString(fmt.Sprintf("- **New behavior:** %s (conf=%.2f)\n", o.NewIntent, o.NewConf))
		b.WriteString(fmt.Sprintf("- **Confidence gap:** %.2f\n", math.Abs(o.NewConf-o.OldConf)))
		b.WriteString(fmt.Sprintf("- **Observations:** %d\n\n", o.Count))
	}

	if r.AlertCount > 0 {
		b.WriteString("---\n")
		b.WriteString("_Drift alerts are observational. They do not block or gate anything._\n")
	}
	return b.String()
}

// SuggestAction returns a suggested response to a drift observation.
func (o Observation) SuggestAction() string {
	switch {
	case o.Severity == SeverityNone:
		return "no action needed"
	case o.Severity == SeverityLow:
		return "monitor — may self-correct"
	case o.OldIntent != o.NewIntent && o.NewIntent != "":
		return fmt.Sprintf("user preference may have shifted from %q to %q — consider updating preference store", o.OldIntent, o.NewIntent)
	case o.NewConf < o.OldConf:
		return "confidence declining — check if skill quality has degraded"
	default:
		return "review calibration data for root cause"
	}
}
