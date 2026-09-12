package agent

import (
	"time"

	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
)

type StageStatus string

type RunStatus string

const (
	StagePending StageStatus = "pending"
	StageRunning StageStatus = "running"
	StageSuccess StageStatus = "success"
	StageFailed  StageStatus = "failed"

	RunQueued    RunStatus = "queued"
	RunRunning   RunStatus = "running"
	RunSucceeded RunStatus = "succeeded"
	RunFailed    RunStatus = "failed"
)

type TraceStep struct {
	Name      string      `json:"name"`
	Status    StageStatus `json:"status"`
	StartedAt time.Time   `json:"started_at,omitempty"`
	EndedAt   time.Time   `json:"ended_at,omitempty"`
	Error     string      `json:"error,omitempty"`
	Meta      any         `json:"meta,omitempty"`
}

type RunResult struct {
	SearchResults []search.Result `json:"search_results"`
	Synthesis     synth.Result    `json:"synthesis"`
}

type RunRecord struct {
	ID        string      `json:"id"`
	Query     string      `json:"query"`
	Status    RunStatus   `json:"status"`
	Result    *RunResult  `json:"result,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Trace     []TraceStep `json:"trace"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}
