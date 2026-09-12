package agent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"sync"
	"time"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
	"github.com/Priyank911/membrow/apps/api/internal/storage"
)

type Service struct {
	store      storage.RunStore
	workflow   *Orchestrator
	runTimeout time.Duration
	logger     *slog.Logger
	baseCtx    context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewService(store storage.RunStore, workflow *Orchestrator, runTimeout time.Duration, logger *slog.Logger) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{store: store, workflow: workflow, runTimeout: runTimeout, logger: logger, baseCtx: ctx, cancel: cancel}
}

func (s *Service) StartRun(ctx context.Context, query string, limit int) (string, error) {
	if query == "" {
		return "", apperrors.InvalidInput("query is required", map[string]string{"field": "query"})
	}
	runID := newRunID()
	now := time.Now().UTC()
	record := agentToStorageRun(RunRecord{ID: runID, Query: query, Status: RunRunning, Trace: []TraceStep{}, CreatedAt: now, UpdatedAt: now})
	if err := s.store.CreateRun(ctx, record); err != nil {
		return "", apperrors.Internal("failed to persist run", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		runCtx, cancel := context.WithTimeout(s.baseCtx, s.runTimeout)
		defer cancel()

		result, trace, err := s.workflow.Run(runCtx, query, limit)
		storageTrace := make([]storage.TraceStep, 0, len(trace))
		for _, item := range trace {
			storageTrace = append(storageTrace, storage.TraceStep{
				Name:      item.Name,
				Status:    storage.StageStatus(item.Status),
				StartedAt: item.StartedAt,
				EndedAt:   item.EndedAt,
				Error:     item.Error,
				Meta:      item.Meta,
			})
		}
		if err != nil {
			appErr := apperrors.AsAppError(err)
			failure := &storage.ErrorInfo{Code: string(appErr.Code), Message: appErr.Message, Retryable: appErr.Retryable}
			if updateErr := s.store.CompleteRun(context.Background(), runID, storage.RunStatusFailed, nil, storageTrace, failure); updateErr != nil {
				s.logger.Error("failed to persist failed run", "run_id", runID, "error", updateErr)
			}
			return
		}
		storageResult := &storage.RunResult{
			SearchResults: result.SearchResults,
			Synthesis:     result.Synthesis,
		}
		if updateErr := s.store.CompleteRun(context.Background(), runID, storage.RunStatusSucceeded, storageResult, storageTrace, nil); updateErr != nil {
			s.logger.Error("failed to persist completed run", "run_id", runID, "error", updateErr)
		}
	}()

	return runID, nil
}

func (s *Service) GetRun(ctx context.Context, runID string) (*RunRecord, error) {
	record, err := s.store.GetRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	agentRecord := storageToAgentRun(record)
	return &agentRecord, nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.cancel()
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func newRunID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().UTC().Format("20060102150405.000")
	}
	return hex.EncodeToString(buf)
}

func agentToStorageRun(in RunRecord) storage.RunRecord {
	var result *storage.RunResult
	if in.Result != nil {
		result = &storage.RunResult{
			SearchResults: in.Result.SearchResults,
			Synthesis:     in.Result.Synthesis,
		}
	}
	var runErr *storage.ErrorInfo
	if in.Error != nil {
		runErr = &storage.ErrorInfo{
			Code:      in.Error.Code,
			Message:   in.Error.Message,
			Retryable: in.Error.Retryable,
		}
	}
	trace := make([]storage.TraceStep, 0, len(in.Trace))
	for _, item := range in.Trace {
		trace = append(trace, storage.TraceStep{
			Name:      item.Name,
			Status:    storage.StageStatus(item.Status),
			StartedAt: item.StartedAt,
			EndedAt:   item.EndedAt,
			Error:     item.Error,
			Meta:      item.Meta,
		})
	}
	return storage.RunRecord{
		ID: in.ID,
		Query: in.Query,
		Status: storage.RunStatus(in.Status),
		Result: result,
		Error: runErr,
		Trace: trace,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func storageToAgentRun(in storage.RunRecord) RunRecord {
	var result *RunResult
	if in.Result != nil {
		result = &RunResult{
			SearchResults: in.Result.SearchResults,
			Synthesis:     in.Result.Synthesis,
		}
	}
	var runErr *ErrorInfo
	if in.Error != nil {
		runErr = &ErrorInfo{
			Code:      in.Error.Code,
			Message:   in.Error.Message,
			Retryable: in.Error.Retryable,
		}
	}
	trace := make([]TraceStep, 0, len(in.Trace))
	for _, item := range in.Trace {
		trace = append(trace, TraceStep{
			Name:      item.Name,
			Status:    StageStatus(item.Status),
			StartedAt: item.StartedAt,
			EndedAt:   item.EndedAt,
			Error:     item.Error,
			Meta:      item.Meta,
		})
	}
	return RunRecord{
		ID: in.ID,
		Query: in.Query,
		Status: RunStatus(in.Status),
		Result: result,
		Error: runErr,
		Trace: trace,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}
