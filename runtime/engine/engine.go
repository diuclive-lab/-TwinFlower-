// Package engine is the minimal runtime loop.
package engine

import (
	"context"
	"fmt"
	"strings"

	"twinflower/root/cognition"
	"twinflower/root/providers"
	"twinflower/vascular/skills/tool_selection"
)

// Tool is the interface every tool must implement.
type Tool interface {
	Name() string
	Run(ctx context.Context, args map[string]any) (string, error)
}

// Engine is the minimal runtime.
type Engine struct {
	provider    providers.Provider
	cognitive   *cognition.Profile // model cognitive profile
	tools       map[string]Tool
	skills      map[string]*tool_selection.Skill
}

// New creates an engine. If profile is nil, uses QwenDense default.
func New(p providers.Provider, cp *cognition.Profile) *Engine {
	if cp == nil {
		cp = cognition.QwenDense()
	}
	return &Engine{
		provider:  p,
		cognitive: cp,
		tools:     make(map[string]Tool),
		skills:    make(map[string]*tool_selection.Skill),
	}
}

func (e *Engine) RegisterTool(t Tool)  { e.tools[t.Name()] = t }
func (e *Engine) RegisterSkill(s *tool_selection.Skill) { e.skills[s.Contract().Name] = s }

// Handle processes a user input through the cognitive loop.
func (e *Engine) Handle(ctx context.Context, input string) (string, error) {
	skill := e.selectSkill(input)
	if skill == nil {
		skill = tool_selection.General()
	}

	// Build tool descriptions
	var toolDefs []providers.ToolDef
	for _, name := range skill.Contract().AllowedTools {
		if t, ok := e.tools[name]; ok {
			toolDefs = append(toolDefs, providers.ToolDef{
				Name:        t.Name(),
				Description: t.Name() + " tool",
			})
		}
	}

	// Use cognitive profile to build the intent prompt
	intentPrompt := buildPlanPrompt(e.cognitive, input, toolDefs, skill)
	plan, err := e.provider.Plan(ctx, intentPrompt, toolDefs)
	if err != nil {
		return "", fmt.Errorf("plan failed: %w", err)
	}

	// Check for clarification
	if plan.Content != "" && e.needsClarify(plan) {
		return plan.Content, nil // the model's response IS the clarify question
	}

	// Check forbidden tools
	if plan.Tool != "" && skill.Forbidden(plan.Tool) {
		return "", fmt.Errorf("skill %s forbids tool %q", skill.Contract().Name, plan.Tool)
	}

	// Execute tool
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

	// Finalize using cognitive profile
	finalPrompt := buildFinalPrompt(e.cognitive, input, toolResult)
	response, err := e.provider.Finalize(ctx, finalPrompt, toolResult)
	if err != nil {
		return "", fmt.Errorf("finalize failed: %w", err)
	}

	return response, nil
}

// needsClarify returns true if the plan indicates low confidence.
func (e *Engine) needsClarify(plan *providers.PlanResult) bool {
	if plan.Confidence < e.cognitive.ClarifyThreshold {
		return true
	}
	if plan.Tool == "clarify" {
		return true
	}
	return false
}

// buildPlanPrompt uses the cognitive profile's prompt template.
func buildPlanPrompt(cp *cognition.Profile, input string, tools []providers.ToolDef, skill *tool_selection.Skill) string {
	if cp.PromptTemplates.Plan != "" {
		allowed := make([]string, len(tools))
		for i, t := range tools {
			allowed[i] = t.Name
		}
		toolNames := strings.Join(allowed, ", ")
		// If skill has constraints, add them
		constraint := ""
		if skill != nil {
			constraint = skill.ConstraintPrompt()
		}
		return fmt.Sprintf(cp.PromptTemplates.Plan, toolNames, input) + "\n" + constraint
	}
	return input
}

// buildFinalPrompt uses the cognitive profile's final template.
func buildFinalPrompt(cp *cognition.Profile, input, toolResult string) string {
	if cp.PromptTemplates.Final != "" {
		return fmt.Sprintf(cp.PromptTemplates.Final, input, toolResult)
	}
	return fmt.Sprintf("User: %s\nTool result: %s", input, toolResult)
}

// selectSkill picks the skill based on input keywords (temporary, will use model later).
func (e *Engine) selectSkill(input string) *tool_selection.Skill {
	lower := strings.ToLower(input)
	if strings.Contains(lower, "比较") || strings.Contains(lower, "compare") ||
		strings.Contains(lower, "diff") || strings.Contains(lower, "对比") {
		return tool_selection.FileCompare()
	}
	if strings.Contains(lower, "天气") || strings.Contains(lower, "weather") ||
		strings.Contains(lower, "翻译") || strings.Contains(lower, "translate") ||
		strings.Contains(lower, "股价") || strings.Contains(lower, "stock") ||
		strings.Contains(lower, "汇率") || strings.Contains(lower, "currency") {
		return tool_selection.BusinessQuery()
	}
	return tool_selection.General()
}
