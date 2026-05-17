// Package providers abstracts model access (local or API).
package providers

import "context"

// Provider is the model interface.
type Provider interface {
	Plan(ctx context.Context, prompt string, tools []ToolDef) (*PlanResult, error)
	Finalize(ctx context.Context, prompt string, result string) (string, error)
}

// ToolDef describes a tool to the model.
type ToolDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"`
}

// IntentAlternative is a secondary intent candidate.
type IntentAlternative struct {
	Intent string  `json:"intent"`
	Score  float64 `json:"score"`
}

// PlanResult is what the model decides to do.
type PlanResult struct {
	Tool         string              `json:"tool"`     // selected tool, or "clarify"
	Args         map[string]any      `json:"args"`
	Content      string              `json:"content,omitempty"`       // direct reply or clarify question
	Confidence   float64             `json:"confidence"`
	Alternatives []IntentAlternative `json:"alternatives,omitempty"`  // other candidates
}
