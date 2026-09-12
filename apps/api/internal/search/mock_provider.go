package search

import (
	"context"
	"fmt"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Search(_ context.Context, query string, limit int) ([]Result, error) {
	results := []Result{
		{Title: "Membrow Overview", URL: "https://example.com/membrow-overview", Snippet: "Membrow is an agentic search scaffold with pluggable tools."},
		{Title: "Agentic Workflow Design", URL: "https://example.com/agentic-workflow", Snippet: "Plan, search, fetch, extract, synthesize with traces."},
		{Title: "Robust Go APIs", URL: "https://example.com/go-api", Snippet: "Patterns for typed errors, retries, and observability in Go."},
	}
	if query != "" {
		results = append(results, Result{Title: fmt.Sprintf("Query insight: %s", query), URL: "https://example.com/query", Snippet: "Synthetic mock result for local development without API keys."})
	}
	if limit <= 0 || limit > len(results) {
		limit = len(results)
	}
	return results[:limit], nil
}
