// Package tool_selection decides which tool to use based on intent + constraints.
package tool_selection

import "strings"

// Contract defines what a skill allows and forbids.
type Contract struct {
	Name           string
	Goal           string
	AllowedTools   []string
	ForbiddenTools []string
	PromptHint     string // injected into model prompt
}

// Skill is a tool orchestration contract.
// It doesn't call the model — it constrains the model's choices.
type Skill struct {
	contract Contract
}

func New(c Contract) *Skill {
	return &Skill{contract: c}
}

func (s *Skill) Contract() Contract { return s.contract }

// Allowed returns true if the tool is permitted by this contract.
func (s *Skill) Allowed(tool string) bool {
	for _, a := range s.contract.AllowedTools {
		if a == tool {
			return true
		}
	}
	// If no allowed list specified, everything is allowed
	if len(s.contract.AllowedTools) == 0 {
		return !s.Forbidden(tool)
	}
	return false
}

// Forbidden returns true if the tool is explicitly prohibited.
func (s *Skill) Forbidden(tool string) bool {
	for _, f := range s.contract.ForbiddenTools {
		if f == tool {
			return true
		}
	}
	return false
}

// ConstraintPrompt returns a string describing what this skill allows/forbids.
func (s *Skill) ConstraintPrompt() string {
	var b strings.Builder
	b.WriteString("Goal: " + s.contract.Goal + "\n")
	if len(s.contract.AllowedTools) > 0 {
		b.WriteString("Allowed tools: " + strings.Join(s.contract.AllowedTools, ", ") + "\n")
	}
	if len(s.contract.ForbiddenTools) > 0 {
		b.WriteString("Forbidden tools: " + strings.Join(s.contract.ForbiddenTools, ", ") + "\n")
	}
	return b.String()
}

// ── Prebuilt skills ─────────────────────────────────────────────────────

func BusinessQuery() *Skill {
	return New(Contract{
		Name:         "business_query",
		Goal:         "Answer business queries like weather, translation, stock prices, currency conversion",
		AllowedTools: []string{"weather", "translate", "stock", "currency"},
		PromptHint:   "Use the appropriate business tool for the query",
	})
}

func FileCompare() *Skill {
	return New(Contract{
		Name:           "file_compare",
		Goal:           "Compare two files and output differences",
		AllowedTools:   []string{"filesystem_read", "filesystem_search"},
		ForbiddenTools: []string{"browser"},
		PromptHint:     "Only use filesystem tools to read files, never use browser for file comparison",
	})
}

func General() *Skill {
	return New(Contract{
		Name: "general",
		Goal: "Respond to general queries without specialized tools",
	})
}
