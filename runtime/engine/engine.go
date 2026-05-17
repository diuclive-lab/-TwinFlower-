// Package engine is the minimal runtime loop.
package engine

import (
	"context"
	"fmt"
	"strings"

	"twinflower/root/cognition"
	"twinflower/root/cognition/calibration"
	"twinflower/root/cognition/preferences"
	"twinflower/root/providers"
	"twinflower/vascular/skills/clarify"
	"twinflower/vascular/skills/tool_selection"
)

// Tool is the interface every tool must implement.
type Tool interface {
	Name() string
	Run(ctx context.Context, args map[string]any) (string, error)
}

// Engine is the minimal runtime.
type Engine struct {
	provider        providers.Provider
	cognitive       *cognition.Profile
	tools           map[string]Tool
	skills          map[string]*tool_selection.Skill
	prefs           *preferences.Store
}

// SetPreferences attaches a preference store for adaptive clarification.
func (e *Engine) SetPreferences(p *preferences.Store) { e.prefs = p }

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

	var softMsg string // soft execute hint

	// Preference resolution: if model chose no tool, check preferences
	if plan.Tool == "" && e.prefs != nil {
		preferred, adjustedConf := e.prefs.Lookup(input)
		if preferred != "" && adjustedConf >= e.cognitive.ClarifyThreshold {
			plan.Tool = preferred
			plan.Confidence = adjustedConf
			if preferred == "weather" && plan.Args == nil {
				plan.Args = map[string]any{"location": extractCity(input)}
			}
			softMsg = fmt.Sprintf("> 猜你是想查%s%s\n（如果不是请纠正）\n", describeIntent(preferred), extractCity(input))
			calibration.Log(calibration.Record{
				Model:      e.cognitive.Name,
				Intent:     "preference_resolve",
				Tool:       preferred,
				Confidence: adjustedConf,
				Clarified:  false,
				Success:    true,
			})
			goto executeTool
		}
	}

	// Check for clarification
	if plan.Content != "" && e.needsClarify(plan) {
		candidates := extractCandidates(plan)
		question := clarify.BuildQuestion(input, candidates, nil)
		calibration.Log(calibration.Record{
			Model:      e.cognitive.Name,
			Intent:     "clarify",
			Tool:       "clarify",
			Confidence: plan.Confidence,
			Clarified:  true,
			Success:    true,
		})
		return question, nil
	}

executeTool:

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

	// Record calibration data
	calibration.Log(calibration.Record{
		Model:      e.cognitive.Name,
		Intent:     skill.Contract().Name,
		Tool:       plan.Tool,
		Confidence: plan.Confidence,
		Success:    err == nil && toolResult != "",
	})

	return softMsg + response, nil
}

// needsClarify returns true if the plan indicates low confidence or ambiguity.
func (e *Engine) needsClarify(plan *providers.PlanResult) bool {
	// Explicit tool "clarify" selected (model output)
	if plan.Tool == "clarify" {
		return true
	}
	// Low confidence regardless of tool selection
	if plan.Confidence < e.cognitive.ClarifyThreshold {
		return true
	}
	// Check for close alternatives (gap < 0.15 between primary and first alternative)
	if len(plan.Alternatives) > 0 && plan.Tool != plan.Alternatives[0].Intent {
		gap := plan.Confidence - plan.Alternatives[0].Score
		if gap >= 0 && gap < 0.15 {
			return true
		}
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
	if strings.Contains(lower, "搜索") || strings.Contains(lower, "search") ||
		strings.Contains(lower, "查一下") || strings.Contains(lower, "找一下") {
		return tool_selection.General()
	}
	return tool_selection.General()
}

// extractCandidates returns intent candidates from the plan result.
func extractCity(input string) string {
	knownCities := []string{"北京", "上海", "广州", "深圳", "杭州", "成都", "南京", "武汉", "天津", "重庆", "苏州", "西安", "长沙", "郑州", "东莞", "青岛", "沈阳", "宁波", "昆明", "大连"}
	for _, city := range knownCities {
		if strings.Contains(input, city) {
			return city
		}
	}
	return "Beijing"
}

func describeIntent(intent string) string {
	m := map[string]string{
		"weather":         "天气（",
		"search":          "信息（",
		"filesystem_list": "目录（",
		"translate":       "翻译（",
	}
	if v, ok := m[intent]; ok {
		return v
	}
	return intent + "（"
}

func extractCandidates(plan *providers.PlanResult) []string {
	candidates := []string{}
	if plan.Tool != "" && plan.Tool != "clarify" {
		candidates = append(candidates, plan.Tool)
	}
	for _, alt := range plan.Alternatives {
		candidates = append(candidates, alt.Intent)
	}
	if len(candidates) == 0 {
		return []string{"chat", "search"}
	}
	return candidates
}
