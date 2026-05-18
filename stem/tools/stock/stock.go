// Package stock provides a stock price lookup tool using Yahoo Finance (no API key).
package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ── Types ────────────────────────────────────────────────────────────────────

// Request is the input for a stock price lookup.
type Request struct {
	Symbol string `json:"symbol"`
}

// Result is the output of a stock price lookup.
type Result struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePct     float64 `json:"change_pct"`
	PreviousClose float64 `json:"previous_close"`
	Currency      string  `json:"currency"`
}

// Provider abstracts stock price retrieval.
type Provider interface {
	Quote(ctx context.Context, req Request) (Result, error)
}

// ── Tool ─────────────────────────────────────────────────────────────────────

// Tool provides stock price lookups.
type Tool struct {
	Provider Provider
}

// New creates a stock tool with the HTTP provider.
func New() *Tool {
	return &Tool{Provider: &HTTPProvider{}}
}

func (t *Tool) Name() string { return "stock" }

func (t *Tool) Run(ctx context.Context, args map[string]any) (string, error) {
	symbol := ""
	for _, key := range []string{"symbol", "ticker", "stock", "code"} {
		if v, ok := args[key]; ok {
			symbol, _ = v.(string)
			break
		}
	}
	if symbol == "" {
		return "", fmt.Errorf("stock: symbol is required")
	}

	result, err := t.Provider.Quote(ctx, Request{Symbol: symbol})
	if err != nil {
		return "", fmt.Errorf("stock: %w", err)
	}

	changeSign := "+"
	if result.Change < 0 {
		changeSign = ""
	}
	return fmt.Sprintf("%s (%s): %.2f %s (%s%.2f, %s%.2f%%)",
		result.Name, result.Symbol, result.Price, result.Currency,
		changeSign, result.Change, changeSign, result.ChangePct), nil
}

// ── HTTP Provider ────────────────────────────────────────────────────────────

// HTTPProvider fetches stock data from Yahoo Finance.
type HTTPProvider struct{}

func (p *HTTPProvider) Quote(ctx context.Context, req Request) (Result, error) {
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", req.Symbol)
	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return Result{}, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(httpReq)
	if err != nil {
		return Result{}, fmt.Errorf("yahoo finance: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return Result{}, fmt.Errorf("yahoo finance returned %d for %s", resp.StatusCode, req.Symbol)
	}

	var chartResp struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Symbol             string  `json:"symbol"`
					RegularMarketPrice float64 `json:"regularMarketPrice"`
					PreviousClose      float64 `json:"previousClose"`
					Currency           string  `json:"currency"`
					ShortName          string  `json:"shortName"`
				} `json:"meta"`
			} `json:"result"`
			Error any `json:"error"`
		} `json:"chart"`
	}
	if err := json.Unmarshal(body, &chartResp); err != nil {
		return Result{}, fmt.Errorf("parse response: %w", err)
	}

	if len(chartResp.Chart.Result) == 0 {
		return Result{}, fmt.Errorf("no data for symbol: %s", req.Symbol)
	}

	meta := chartResp.Chart.Result[0].Meta
	price := meta.RegularMarketPrice
	prevClose := meta.PreviousClose
	change := price - prevClose
	changePct := 0.0
	if prevClose > 0 {
		changePct = (change / prevClose) * 100
	}

	name := meta.ShortName
	if name == "" {
		name = req.Symbol
	}

	return Result{
		Symbol: req.Symbol, Name: name, Price: price,
		Change: change, ChangePct: changePct, PreviousClose: prevClose, Currency: meta.Currency,
	}, nil
}
