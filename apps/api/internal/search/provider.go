package search

import (
	"context"
)

type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type Provider interface {
	Search(ctx context.Context, query string, limit int) ([]Result, error)
}
