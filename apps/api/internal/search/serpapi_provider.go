package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
)

type SerpAPIProvider struct {
	apiKey string
	client *http.Client
}

type serpAPIResponse struct {
	OrganicResults []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"organic_results"`
}

func NewSerpAPIProvider(apiKey string) *SerpAPIProvider {
	return &SerpAPIProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 8 * time.Second},
	}
}

func (p *SerpAPIProvider) Search(ctx context.Context, query string, limit int) ([]Result, error) {
	if p.apiKey == "" {
		return nil, apperrors.DependencyFailed("search provider key is missing", nil, false)
	}
	if limit <= 0 {
		limit = 5
	}
	endpoint := fmt.Sprintf("https://serpapi.com/search.json?engine=google&q=%s&num=%d&api_key=%s", url.QueryEscape(query), limit, url.QueryEscape(p.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, apperrors.Internal("failed to create search request", err)
	}
	res, err := p.client.Do(req)
	if err != nil {
		return nil, apperrors.DependencyFailed("search provider request failed", err, true)
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		return nil, apperrors.DependencyFailed(fmt.Sprintf("search provider returned status %d", res.StatusCode), nil, res.StatusCode >= 500)
	}

	var payload serpAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, apperrors.DependencyFailed("failed to decode search provider response", err, false)
	}

	results := make([]Result, 0, len(payload.OrganicResults))
	for _, item := range payload.OrganicResults {
		results = append(results, Result{Title: item.Title, URL: item.Link, Snippet: item.Snippet})
	}
	if len(results) == 0 {
		return []Result{}, nil
	}
	if limit > len(results) {
		limit = len(results)
	}
	return results[:limit], nil
}
