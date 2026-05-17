// Package search provides a web search tool using DuckDuckGo (no API key).
package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const duckduckgoURL = "https://html.duckduckgo.com/html/"

type Tool struct{}

func New() *Tool { return &Tool{} }

func (t *Tool) Name() string { return "search" }

func (t *Tool) Run(ctx context.Context, args map[string]any) (string, error) {
	query := ""
	for _, key := range []string{"query", "q", "search", "关键词", "搜索词"} {
		if v, ok := args[key]; ok {
			query, _ = v.(string)
			break
		}
	}
	if query == "" {
		return "", fmt.Errorf("search: query is required")
	}

	v := url.Values{}
	v.Set("q", query)
	resp, err := http.PostForm(duckduckgoURL, v)
	if err != nil {
		return "", fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	text := extractText(string(body))
	if len(text) > 2000 {
		text = text[:2000] + "..."
	}
	return text, nil
}

func extractText(html string) string {
	var results []string
	lines := strings.Split(html, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Extract result snippets from DDG HTML
		if strings.Contains(trimmed, "class=\"result__snippet\"") ||
			strings.Contains(trimmed, "class=\"result__title\"") {
			text := stripTags(trimmed)
			if text != "" {
				results = append(results, text)
			}
		}
	}
	if len(results) == 0 {
		return "No search results found"
	}
	return strings.Join(results, "\n")
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, c := range s {
		if c == '<' {
			inTag = true
			continue
		}
		if c == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(c)
		}
	}
	return strings.TrimSpace(b.String())
}
