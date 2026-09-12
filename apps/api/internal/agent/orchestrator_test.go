package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Priyank911/membrow/apps/api/internal/extract"
	"github.com/Priyank911/membrow/apps/api/internal/fetch"
	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
)

type fakeSearch struct{}

type fakeFetcher struct{}

type fakeExtractor struct{}

type fakeSynth struct{}

type failingSearch struct{}

func (f *fakeSearch) Search(_ context.Context, _ string, _ int) ([]search.Result, error) {
	return []search.Result{{Title: "t1", URL: "https://example.com", Snippet: "snippet"}}, nil
}

func (f *fakeFetcher) Fetch(_ context.Context, url string) (fetch.Document, error) {
	return fetch.Document{URL: url, Content: "<h1>Hello</h1> world"}, nil
}

func (f *fakeExtractor) Extract(_ context.Context, d fetch.Document) (extract.Output, error) {
	return extract.Output{URL: d.URL, Snippet: "hello world", Text: "hello world"}, nil
}

func (f *fakeSynth) Synthesize(_ context.Context, _ string, items []extract.Output) (synth.Result, error) {
	return synth.Result{Answer: "ok", Sources: items}, nil
}

func (f *failingSearch) Search(_ context.Context, _ string, _ int) ([]search.Result, error) {
	return nil, errors.New("search failed")
}

func TestOrchestratorRunSuccess(t *testing.T) {
	orch := NewOrchestrator(&fakeSearch{}, &fakeFetcher{}, &fakeExtractor{}, &fakeSynth{}, 2, 1, 5*time.Millisecond)

	result, trace, err := orch.Run(context.Background(), "test", 1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || result.Synthesis.Answer == "" {
		t.Fatal("expected synthesis result")
	}
	if len(trace) != 4 {
		t.Fatalf("expected 4 trace steps, got %d", len(trace))
	}
}

func TestOrchestratorRunFailure(t *testing.T) {
	orch := NewOrchestrator(&failingSearch{}, &fakeFetcher{}, &fakeExtractor{}, &fakeSynth{}, 2, 1, 5*time.Millisecond)

	_, trace, err := orch.Run(context.Background(), "test", 1)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(trace) < 2 {
		t.Fatalf("expected trace with failed step, got %d", len(trace))
	}
}
