package agent

import (
	"context"
	"sync"
	"time"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
	"github.com/Priyank911/membrow/apps/api/internal/extract"
	"github.com/Priyank911/membrow/apps/api/internal/fetch"
	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
	"github.com/Priyank911/membrow/apps/api/internal/tools"
)

type Orchestrator struct {
	search         search.Provider
	fetcher        fetch.Fetcher
	extractor      extract.Extractor
	synthesizer    synth.Synthesizer
	workerPoolSize int
	maxRetries     int
	retryBackoff   time.Duration
}

func NewOrchestrator(
	searchProvider search.Provider,
	fetcher fetch.Fetcher,
	extractor extract.Extractor,
	synthesizer synth.Synthesizer,
	workerPoolSize int,
	maxRetries int,
	retryBackoff time.Duration,
) *Orchestrator {
	if workerPoolSize < 1 {
		workerPoolSize = 1
	}
	return &Orchestrator{
		search:         searchProvider,
		fetcher:        fetcher,
		extractor:      extractor,
		synthesizer:    synthesizer,
		workerPoolSize: workerPoolSize,
		maxRetries:     maxRetries,
		retryBackoff:   retryBackoff,
	}
}

func (o *Orchestrator) Run(ctx context.Context, query string, limit int) (*RunResult, []TraceStep, error) {
	trace := make([]TraceStep, 0, 5)
	planStep := startStep("plan")
	planStep.Meta = map[string]any{"query": query, "limit": limit}
	trace = append(trace, finishStep(planStep, nil))

	searchStep := startStep("search")
	searchResults, err := tools.WithRetry(ctx, o.maxRetries+1, o.retryBackoff, func(runCtx context.Context) ([]search.Result, error) {
		return o.search.Search(runCtx, query, limit)
	})
	trace = append(trace, finishStep(searchStep, err))
	if err != nil {
		return nil, trace, err
	}

	fetchStep := startStep("fetch")
	documents, err := o.fetchAndExtract(ctx, searchResults)
	trace = append(trace, finishStep(fetchStep, err))
	if err != nil {
		return nil, trace, err
	}

	synthStep := startStep("synthesize")
	synthesis, err := tools.WithRetry(ctx, o.maxRetries+1, o.retryBackoff, func(runCtx context.Context) (synth.Result, error) {
		return o.synthesizer.Synthesize(runCtx, query, documents)
	})
	trace = append(trace, finishStep(synthStep, err))
	if err != nil {
		return nil, trace, err
	}

	return &RunResult{SearchResults: searchResults, Synthesis: synthesis}, trace, nil
}

func (o *Orchestrator) fetchAndExtract(ctx context.Context, results []search.Result) ([]extract.Output, error) {
	if len(results) == 0 {
		return []extract.Output{}, nil
	}
	type jobResult struct {
		value extract.Output
		err   error
	}

	resultsCh := make(chan jobResult, len(results))
	sem := make(chan struct{}, o.workerPoolSize)
	var wg sync.WaitGroup

	for _, item := range results {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				resultsCh <- jobResult{err: apperrors.Timeout("fetch stage canceled", ctx.Err())}
				return
			}
			defer func() { <-sem }()

			doc, err := tools.WithRetry(ctx, o.maxRetries+1, o.retryBackoff, func(runCtx context.Context) (fetch.Document, error) {
				return o.fetcher.Fetch(runCtx, item.URL)
			})
			if err != nil {
				resultsCh <- jobResult{err: err}
				return
			}
			out, err := tools.WithRetry(ctx, o.maxRetries+1, o.retryBackoff, func(runCtx context.Context) (extract.Output, error) {
				return o.extractor.Extract(runCtx, doc)
			})
			if err != nil {
				resultsCh <- jobResult{err: err}
				return
			}
			out.Title = item.Title
			if out.Snippet == "" {
				out.Snippet = item.Snippet
			}
			resultsCh <- jobResult{value: out}
		}()
	}

	wg.Wait()
	close(resultsCh)

	outputs := make([]extract.Output, 0, len(results))
	for result := range resultsCh {
		if result.err != nil {
			return nil, result.err
		}
		outputs = append(outputs, result.value)
	}
	return outputs, nil
}

func startStep(name string) TraceStep {
	return TraceStep{Name: name, Status: StageRunning, StartedAt: time.Now().UTC()}
}

func finishStep(step TraceStep, err error) TraceStep {
	step.EndedAt = time.Now().UTC()
	if err != nil {
		step.Status = StageFailed
		step.Error = err.Error()
		return step
	}
	step.Status = StageSuccess
	return step
}
