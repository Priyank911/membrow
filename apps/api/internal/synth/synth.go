package synth

import (
	"context"
	"fmt"
	"strings"

	"github.com/Priyank911/membrow/apps/api/internal/extract"
)

type Result struct {
	Answer  string            `json:"answer"`
	Sources []extract.Output  `json:"sources"`
}

type Synthesizer interface {
	Synthesize(ctx context.Context, query string, items []extract.Output) (Result, error)
}

type BasicSynth struct{}

func NewBasicSynth() *BasicSynth {
	return &BasicSynth{}
}

func (s *BasicSynth) Synthesize(_ context.Context, query string, items []extract.Output) (Result, error) {
	if len(items) == 0 {
		return Result{Answer: "No sources were available to synthesize a response.", Sources: []extract.Output{}}, nil
	}
	parts := make([]string, 0, len(items))
	for i, item := range items {
		parts = append(parts, fmt.Sprintf("%d. %s", i+1, item.Snippet))
	}
	return Result{
		Answer:  fmt.Sprintf("Query: %s\n\nSummary from collected sources:\n%s", query, strings.Join(parts, "\n")),
		Sources: items,
	}, nil
}
