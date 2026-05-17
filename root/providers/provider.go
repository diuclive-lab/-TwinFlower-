// Package providers abstracts model access (local or API).
package providers

import "context"

// Provider is the model interface. The only thing a model does: plan and finalize.
type Provider interface {
	// Plan sends messages and expects a structured tool call.
	Plan(ctx context.Context, prompt string, tools []ToolDef) (*PlanResult, error)

	// Finalize generates the final response after tool execution.
	Finalize(ctx context.Context, prompt string, result string) (string, error)
}

// ToolDef describes a tool to the model.
type ToolDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"`
}

// PlanResult is what the model decides to do.
type PlanResult struct {
	Tool     string         `json:"tool"`
	Args     map[string]any `json:"args"`
	Content  string         `json:"content,omitempty"` // direct reply if no tool
	Confidence float64      `json:"confidence"`
}
