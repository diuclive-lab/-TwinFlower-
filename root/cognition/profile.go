// Package cognition adapts the system's behavior to different model architectures.
//
// Dense models (Qwen, Gemma dense): need structured, explicit reasoning prompts.
// MoE models (DeepSeek, Mixtral): need gated decomposition, higher clarify threshold.
//
// This is not about how to call the model (that's providers).
// This is about how to talk to the model.
package cognition

// Profile describes a model's cognitive characteristics.
// It determines how the system constructs prompts, sets thresholds, and interprets responses.
type Profile struct {
	Name             string   `json:"name"`
	Architecture     string   `json:"architecture"`     // "dense" or "moe"
	Strengths        []string `json:"strengths"`         // e.g. "coding", "chinese", "deterministic"
	Weaknesses       []string `json:"weaknesses"`        // e.g. "ambiguity", "over-routing"
	ClarifyThreshold float64  `json:"clarify_threshold"` // lower = more aggressive clarification
	ToolBias         string   `json:"tool_bias"`         // "conservative", "balanced", "aggressive"
	ChainStyle       string   `json:"chain_style"`       // "explicit", "gated", "free"

	// PromptTemplates defines how different stages talk to this model.
	PromptTemplates PromptSet `json:"prompt_templates"`
}

// PromptSet contains the prompt templates for each stage.
type PromptSet struct {
	Intent  string `json:"intent"`  // intent classification prompt
	Plan    string `json:"plan"`    // tool selection prompt
	Final   string `json:"final"`   // response formatting prompt
	Clarify string `json:"clarify"` // clarification prompt
}

// ── Prebuilt profiles ─────────────────────────────────────────────────────

func QwenDense() *Profile {
	return &Profile{
		Name:             "qwen_dense",
		Architecture:     "dense",
		Strengths:        []string{"coding", "chinese", "deterministic"},
		Weaknesses:       []string{"ambiguity", "implicit intent"},
		ClarifyThreshold: 0.45,
		ToolBias:         "balanced",
		ChainStyle:       "explicit",
		PromptTemplates: PromptSet{
			Intent: `Determine the user's intent from their request.
Output format:
{
  "intent": "tool_name or chat",
  "confidence": 0.0-1.0,
  "evidence": ["keyword1", "keyword2"],
  "alternatives": ["other_possible_intents"]
}
Rules:
- Be precise about what you know.
- If multiple intents are possible, list them as alternatives.
- If uncertainty is high (confidence < %f), set intent to "clarify".
- Available tools: %s`,
			Plan: `You have access to the following tools: %s
Choose the correct tool based on the user's intent.
Constraints:
- Only use tools from the allowed list.
- If no tool matches, respond directly.
User request: %s`,
			Final: `The user asked: %s
The tool returned: %s
Format this as a natural, conversational response.`,
		},
	}
}

func DeepSeekMoE() *Profile {
	return &Profile{
		Name:             "deepseek_moe",
		Architecture:     "moe",
		Strengths:        []string{"mixed reasoning", "broad knowledge", "creative"},
		Weaknesses:       []string{"over-routing", "shortcut inference", "early convergence"},
		ClarifyThreshold: 0.70,
		ToolBias:         "conservative",
		ChainStyle:       "gated",
		PromptTemplates: PromptSet{
			Intent: `First classify the question type: factual / planning / tool-use / ambiguous.
Then respond with:
{
  "intent": "tool_name or chat or clarify",
  "confidence": 0.0-1.0,
  "evidence": ["keyword1"],
  "alternatives": [],
  "question_type": "tool-use"
}
If the request is ambiguous, always set intent to "clarify".
Available tools: %s`,
			Plan: `You selected tool: %s
Your task: execute this tool call precisely.
Do not add extra interpretation.
Tool name: %s
Arguments: determine from user request: %s`,
			Final: `Summarize this result naturally: %s`,
		},
	}
}

// ── Registry ──────────────────────────────────────────────────────────────

// Registry holds all known cognitive profiles.
type Registry struct {
	profiles map[string]*Profile
}

func NewRegistry() *Registry {
	r := &Registry{profiles: make(map[string]*Profile)}
	r.Register(QwenDense())
	r.Register(DeepSeekMoE())
	return r
}

func (r *Registry) Register(p *Profile) {
	r.profiles[p.Name] = p
}

func (r *Registry) Get(name string) *Profile {
	if p, ok := r.profiles[name]; ok {
		return p
	}
	return r.profiles["qwen_dense"] // default
}

// GetByModel returns a profile based on model name pattern matching.
func (r *Registry) GetByModel(modelName string) *Profile {
	// TODO: implement model name → profile mapping
	return r.profiles["qwen_dense"]
}
