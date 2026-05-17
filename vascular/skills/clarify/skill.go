// Package clarify is the formal uncertainty contract.
// When the system doesn't know, it doesn't guess — it asks.
package clarify

import "fmt"

// Request is what the cognition layer produces when confidence is low.
type Request struct {
	NeedsClarify bool     `json:"needs_clarify"`
	Question     string   `json:"question"`
	Candidates   []string `json:"candidates"`
	Confidence   float64  `json:"confidence"`
	Evidence     []string `json:"evidence"`
}

// Response is what the user answers (recorded for calibration).
type Response struct {
	Input     string   `json:"input"`
	Selected  string   `json:"selected"`
	Candidates []string `json:"candidates"`
	Timestamp string   `json:"timestamp"`
}

// BuildQuestion formats a clarify question from candidates and evidence.
func BuildQuestion(input string, candidates []string, evidence []string) string {
	if len(candidates) == 0 {
		return "请更具体地描述您想做什么"
	}
	question := "您是想"
	for i, c := range candidates {
		if i > 0 && i == len(candidates)-1 {
			question += "，还是"
		} else if i > 0 {
			question += "、"
		}
		question += describe(c)
	}
	question += "？"
	return question
}

func describe(candidate string) string {
	m := map[string]string{
		"weather":         "查天气",
		"search":          "搜索信息",
		"filesystem_list": "列出文件",
		"filesystem_read": "读取文件",
		"filesystem_search": "搜索文件",
		"translate":       "翻译",
		"stock":           "查股价",
		"currency":        "查汇率",
		"chat":            "聊天",
	}
	if v, ok := m[candidate]; ok {
		return v
	}
	return candidate
}

// Contract returns the skill contract for clarify.
func Contract() string {
	return fmt.Sprintf(`Clarify skill:
- When the user's intent is ambiguous, ask a clarifying question.
- Output format: {"needs_clarify":true,"question":"...","candidates":["..."],"confidence":0.42,"evidence":["..."]}
- Do NOT guess the intent when uncertain.
- Present 2-3 specific options based on what you detect.`)
}
