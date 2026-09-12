package extract

import (
	"context"
	"regexp"
	"strings"

	"github.com/Priyank911/membrow/apps/api/internal/fetch"
)

type Output struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Text    string `json:"text"`
}

type Extractor interface {
	Extract(ctx context.Context, document fetch.Document) (Output, error)
}

type BasicExtractor struct {
	tagRegex *regexp.Regexp
}

func NewBasicExtractor() *BasicExtractor {
	return &BasicExtractor{tagRegex: regexp.MustCompile(`<[^>]*>`) }
}

func (e *BasicExtractor) Extract(_ context.Context, document fetch.Document) (Output, error) {
	plain := e.tagRegex.ReplaceAllString(document.Content, " ")
	plain = strings.Join(strings.Fields(plain), " ")
	if len(plain) > 2_000 {
		plain = plain[:2_000]
	}
	snippet := plain
	if len(snippet) > 180 {
		snippet = snippet[:180]
	}
	return Output{
		URL: document.URL,
		Title: document.URL,
		Snippet: snippet,
		Text: plain,
	}, nil
}
