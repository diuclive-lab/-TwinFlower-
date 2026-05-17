package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LocalProvider calls a local llama.cpp server via OpenAI-compatible API.
type LocalProvider struct {
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
	APIKey  string `json:"api_key"`
	client  *http.Client
}

func NewLocalProvider(baseURL, model, apiKey string) *LocalProvider {
	return &LocalProvider{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Model:   model,
		APIKey:  apiKey,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *LocalProvider) Plan(ctx context.Context, prompt string, tools []ToolDef) (*PlanResult, error) {
	body := map[string]any{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 256,
		"temperature": 0.3,
	}
	if len(tools) > 0 {
		// Wrap in OpenAI format: {type: "function", function: {...}}
		openaiTools := make([]map[string]any, len(tools))
		for i, t := range tools {
			openaiTools[i] = map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        t.Name,
					"description": t.Description,
					"parameters":  t.Parameters,
				},
			}
		}
		body["tools"] = openaiTools
	}
	return p.call(ctx, body)
}

func (p *LocalProvider) Finalize(ctx context.Context, prompt string, result string) (string, error) {
	body := map[string]any{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "user", "content": fmt.Sprintf("%s\n\nTool result: %s\n\nPlease format this as a natural response.", prompt, result)},
		},
		"max_tokens": 512,
		"temperature": 0.5,
	}
	res, err := p.call(ctx, body)
	if err != nil {
		return "", err
	}
	return res.Content, nil
}

func (p *LocalProvider) call(ctx context.Context, body any) (*PlanResult, error) {
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("llama call failed: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("llama status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					Function struct {
						Name      string          `json:"name"`
						Arguments json.RawMessage `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	msg := result.Choices[0].Message
	pr := &PlanResult{Content: msg.Content, Confidence: 0.8}

	if len(msg.ToolCalls) > 0 {
		tc := msg.ToolCalls[0]
		pr.Tool = tc.Function.Name
		// Arguments can be a JSON object or a stringified JSON string
		raw := []byte(tc.Function.Arguments)
		if err := json.Unmarshal(raw, &pr.Args); err != nil {
			var s string
			if json.Unmarshal(raw, &s) == nil {
				json.Unmarshal([]byte(s), &pr.Args)
			}
		}
	}
	return pr, nil
}
