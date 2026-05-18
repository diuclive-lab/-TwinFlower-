// Package search_skill handles web search with procedural intelligence:
// query extraction -> ambiguity detection -> clarify/soft execute -> search -> result shaping.
package search_skill

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ── Types ────────────────────────────────────────────────────────────────────

// AmbiguityLevel represents how much ambiguity exists in the search query.
type AmbiguityLevel int

const (
	LevelLow    AmbiguityLevel = iota // clear intent, execute directly
	LevelMedium                       // some ambiguity, soft execute with hint
	LevelHigh                         // too ambiguous, must clarify
)

// ToolRunner is the interface a search tool must satisfy (matches engine.Tool).
type ToolRunner interface {
	Name() string
	Run(ctx context.Context, args map[string]any) (string, error)
}

// Skill handles web search requests with procedural intelligence.
type Skill struct {
	tool ToolRunner
}

// New creates a search skill backed by the given search tool.
func New(tool ToolRunner) *Skill {
	return &Skill{tool: tool}
}

// ── Ambiguous Term Registry ──────────────────────────────────────────────────

// ambiguousTerms lists terms that map to multiple possible interpretations.
var ambiguousTerms = map[string][]string{
	"苹果":  {"苹果公司", "苹果手机", "苹果（水果）"},
	"小米":  {"小米公司", "小米（谷物）"},
	"华为":  {"华为公司", "华为手机"},
	"三星":  {"三星公司", "三星手机"},
	"大疆":  {"大疆创新"},
	"特斯拉": {"特斯拉汽车"},
	"阿里巴巴": {"阿里巴巴集团"},
	"腾讯":  {"腾讯公司"},
	"百度":  {"百度公司"},
	"京东":  {"京东电商"},
	"抖音":  {"抖音短视频"},
	"美团":  {"美团"},
	"滴滴":  {"滴滴出行"},
	"联想":  {"联想公司", "联想电脑"},
	"比亚迪": {"比亚迪汽车"},
	"中兴":  {"中兴通讯"},
}

// disambiguators are qualifiers that resolve known ambiguous terms.
var disambiguators = []string{
	"公司", "手机", "股价", "水果", "集团", "汽车", "电商", "谷物",
	"电脑", "通讯", "创新", "短视频", "出行",
	"stock", "price", "company", "phone", "fruit",
}

// ── Handle ───────────────────────────────────────────────────────────────────

// Handle processes a search request through the full pipeline.
func (s *Skill) Handle(ctx context.Context, input string) (string, error) {
	query := extractQuery(input)

	level, candidates := detectAmbiguity(query)

	switch level {
	case LevelHigh:
		return buildClarifyQuestion(query, candidates), nil

	case LevelMedium:
		searchQuery := query
		if len(candidates) > 0 {
			searchQuery = candidates[0]
		}
		result, err := s.executeSearch(ctx, searchQuery)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("> 搜索了 %q（如果不是请纠正）\n%s", searchQuery, result), nil

	default:
		result, err := s.executeSearch(ctx, query)
		if err != nil {
			return "", err
		}
		return result, nil
	}
}

// ── Query Extraction ─────────────────────────────────────────────────────────

// extractQuery pulls the search query from natural language input.
func extractQuery(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	text := strings.TrimSpace(input)

	// Chinese prefixes: if prefix matches but nothing follows, it's an empty query
	cnPrefixes := []string{"搜索", "查一下", "查找", "搜一下", "帮我查", "查询", "找一下", "搜搜"}
	for _, p := range cnPrefixes {
		if idx := strings.Index(text, p); idx >= 0 {
			after := trimPrefix(text, p, idx)
			if after == "" {
				return ""
			}
			return after
		}
	}

	// English prefixes
	enPrefixes := []string{"search for ", "search ", "find ", "look up ", "find about ", "google "}
	for _, p := range enPrefixes {
		if idx := strings.Index(lower, p); idx >= 0 {
			after := strings.TrimSpace(text[idx+len(p):])
			if after == "" {
				return ""
			}
			return after
		}
	}

	// "怎么样" / "如何" pattern: return the subject before it
	if idx := strings.LastIndex(text, "怎么样"); idx > 0 {
		return strings.TrimSpace(text[:idx])
	}
	if idx := strings.LastIndex(text, "如何"); idx > 0 {
		return strings.TrimSpace(text[:idx]) + " 最新"
	}
	if idx := strings.LastIndex(lower, "how about "); idx >= 0 {
		after := strings.TrimSpace(text[idx+10:])
		if after == "" {
			return ""
		}
		return after
	}
	if idx := strings.LastIndex(lower, "what about "); idx >= 0 {
		after := strings.TrimSpace(text[idx+11:])
		if after == "" {
			return ""
		}
		return after
	}

	// No prefix matched — treat whole input as the query
	return text
}

// trimPrefix removes the prefix at byte offset idx and returns the remainder.
func trimPrefix(text, prefix string, idx int) string {
	return strings.TrimSpace(text[idx+len(prefix):])
}

// ── Ambiguity Detection ──────────────────────────────────────────────────────

// detectAmbiguity analyzes a query and returns the ambiguity level and candidates.
func detectAmbiguity(query string) (AmbiguityLevel, []string) {
	query = strings.TrimSpace(query)
	if query == "" {
		return LevelHigh, nil
	}

	// Check for known ambiguous terms
	if candidates, ok := isKnownAmbiguous(query); ok {
		return LevelHigh, candidates
	}

	// Single CJK character is highly ambiguous
	runeCount := utf8.RuneCountInString(query)
	if runeCount <= 1 {
		return LevelHigh, nil
	}

	// Short single-word queries may lack context (Chinese: no spaces, so only check rune count)
	if runeCount <= 2 {
		return LevelMedium, []string{query}
	}

	return LevelLow, nil
}

// isKnownAmbiguous checks if the query contains a known ambiguous term
// without a disambiguating qualifier.
func isKnownAmbiguous(term string) ([]string, bool) {
	for key, candidates := range ambiguousTerms {
		if strings.Contains(term, key) {
			// Check for disambiguating qualifier
			for _, q := range disambiguators {
				if strings.Contains(term, q) {
					return nil, false
				}
			}
			return candidates, true
		}
	}
	return nil, false
}

// ── Clarify ──────────────────────────────────────────────────────────────────

// buildClarifyQuestion generates a clarification question for ambiguous queries.
func buildClarifyQuestion(query string, candidates []string) string {
	if query == "" {
		return "你想搜索什么？请告诉我关键词。"
	}
	if len(candidates) == 0 {
		return fmt.Sprintf("你想搜索 %q 的哪方面信息？请说得更具体一些。", query)
	}

	var b strings.Builder
	b.WriteString("你是想搜索")

	if len(candidates) > 0 {
		for i, c := range candidates {
			if i > 0 && i == len(candidates)-1 {
				b.WriteString("，还是")
			} else if i > 0 {
				b.WriteString("、")
			}
			b.WriteString(c)
		}
	}

	b.WriteString("？")
	return b.String()
}

// ── Search Execution ─────────────────────────────────────────────────────────

// executeSearch runs the query through the search tool and formats results.
func (s *Skill) executeSearch(ctx context.Context, query string) (string, error) {
	result, err := s.tool.Run(ctx, map[string]any{"query": query})
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	return formatResults(query, result), nil
}

// ── Result Shaping ───────────────────────────────────────────────────────────

// formatResults wraps raw search results with a header and structured formatting.
func formatResults(query string, raw string) string {
	if raw == "" || raw == "No search results found" {
		return fmt.Sprintf("未找到与 %q 相关的结果", query)
	}

	lines := strings.Split(raw, "\n")

	var clean []string
	seen := make(map[string]bool)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		clean = append(clean, trimmed)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "🔍 搜索: %s\n\n", query)

	displayed := 0
	for _, line := range clean {
		if displayed >= 8 {
			break
		}
		runes := []rune(line)
		if len(runes) > 150 {
			line = string(runes[:150]) + "…"
		}
		b.WriteString(line)
		b.WriteString("\n")
		displayed++
	}

	if len(clean) > 8 {
		fmt.Fprintf(&b, "\n... 还有 %d 条结果", len(clean)-8)
	}

	return b.String()
}
