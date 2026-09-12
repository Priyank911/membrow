package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
)

type Document struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

type Fetcher interface {
	Fetch(ctx context.Context, url string) (Document, error)
}

type HTTPFetcher struct {
	client    *http.Client
	userAgent string
}

func NewHTTPFetcher(userAgent string) *HTTPFetcher {
	return &HTTPFetcher{
		client:    &http.Client{Timeout: 10 * time.Second},
		userAgent: userAgent,
	}
}

func (f *HTTPFetcher) Fetch(ctx context.Context, rawURL string) (Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Document{}, apperrors.DependencyFailed("failed to create fetch request", err, false)
	}
	req.Header.Set("User-Agent", f.userAgent)
	res, err := f.client.Do(req)
	if err != nil {
		return Document{}, apperrors.DependencyFailed("fetch request failed", err, true)
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		return Document{}, apperrors.DependencyFailed(fmt.Sprintf("fetch returned status %d", res.StatusCode), nil, res.StatusCode >= 500)
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, 256*1024))
	if err != nil {
		return Document{}, apperrors.DependencyFailed("failed to read fetched content", err, true)
	}
	return Document{URL: rawURL, Content: string(body)}, nil
}
