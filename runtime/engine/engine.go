// Package engine is the minimal runtime loop.
package engine

import (
	"context"
	"fmt"
	"strings"

	"twinflower/root/providers"
	"twinflower/vascular/skills/tool_selection"
)

// Tool is the interface every tool must implement.
type Tool interface {
	Name() string
	Run(ctx context.Context, args map[string]any) (string, error)
}

// Engine is the minimal runtime.
// No shadow, no gate, no freeze, no workflow.
type Engine struct {
	provider providers.Provider
	tools    map[string]Tool
	skills   map[string]*tool_selection.Skill
}

func New(p providers.Provider) *Engine {
	return &Engine{
		provider: p,
		tools:    make(map[string]Tool),
		skills:   make(map[string]*tool_selection.Skill),
	}
}

func (e *Engine) RegisterTool(t Tool)     { e.tools[t.Name()] = t }
func (e *Engine) RegisterSkill(s *tool_selection.Skill) { e.skills[s.Contract().Name] = s }

// Handle processes a user input through the minimal loop.
//  1. Select skill based on intent
//  2. Model plans with skill constraints
//  3. Execute tool if model chose one
//  4. Model finalizes
func (e *Engine) Handle(ctx context.Context, input string) (string, error) {
	// 1. Skill selection
	skill := e.selectSkill(input)
	if skill == nil {
		skill = tool_selection.General()
	}

	// 2. Build tool descriptions for the model
	var toolDefs []providers.ToolDef
	for _, name := range skill.Contract().AllowedTools {
		if t, ok := e.tools[name]; ok {
			toolDefs = append(toolDefs, providers.ToolDef{
				Name:        t.Name(),
				Description: t.Name() + " tool",
			})
		}
	}

	// Build prompt with skill constraints
	prompt := input
	if constraint := skill.ConstraintPrompt(); constraint != "" {
		prompt = constraint + "\nUser request: " + input
	}

	// 3. Model plans
	plan, err := e.provider.Plan(ctx, prompt, toolDefs)
	if err != nil {
		return "", fmt.Errorf("plan failed: %w", err)
	}

	// 4. Check forbidden tools
	if plan.Tool != "" && skill.Forbidden(plan.Tool) {
		return "", fmt.Errorf("skill %s forbids tool %q", skill.Contract().Name, plan.Tool)
	}

	// 5. Execute tool
	var toolResult string
	if plan.Tool != "" {
		t, ok := e.tools[plan.Tool]
		if !ok {
			return "", fmt.Errorf("tool %q not registered", plan.Tool)
		}
		toolResult, err = t.Run(ctx, plan.Args)
		if err != nil {
			return "", fmt.Errorf("tool %q failed: %w", plan.Tool, err)
		}
	}

	// 6. Finalize (model formats the final response)
	response, err := e.provider.Finalize(ctx, input, toolResult)
	if err != nil {
		return "", fmt.Errorf("finalize failed: %w", err)
	}

	return response, nil
}

// selectSkill picks the skill based on input keywords (temporary, will use model later).
func (e *Engine) selectSkill(input string) *tool_selection.Skill {
	lower := strings.ToLower(input)
	// File operations → file_compare skill
	if strings.Contains(lower, "比较") || strings.Contains(lower, "compare") ||
		strings.Contains(lower, "diff") || strings.Contains(lower, "对比") {
		return tool_selection.FileCompare()
	}
	// Business queries → business_query skill
	if strings.Contains(lower, "天气") || strings.Contains(lower, "weather") ||
		strings.Contains(lower, "翻译") || strings.Contains(lower, "translate") ||
		strings.Contains(lower, "股价") || strings.Contains(lower, "stock") ||
		strings.Contains(lower, "汇率") || strings.Contains(lower, "currency") {
		return tool_selection.BusinessQuery()
	}
	return tool_selection.General()
}
