package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
)

type RunStatus string

const (
	RunStatusRunning   RunStatus = "running"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
)

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

type RunStore interface {
	CreateRun(ctx context.Context, run RunRecord) error
	CompleteRun(ctx context.Context, runID string, status RunStatus, result *RunResult, trace []TraceStep, runErr *ErrorInfo) error
	GetRun(ctx context.Context, runID string) (RunRecord, error)
	Close() error
}

type StageStatus string

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

type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

type SQLiteRunStore struct {
	db *sql.DB
}

func NewSQLiteRunStore(path string) (*SQLiteRunStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	store := &SQLiteRunStore{db: db}
	if err := store.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteRunStore) migrate() error {
	query := `
CREATE TABLE IF NOT EXISTS runs (
  id TEXT PRIMARY KEY,
  query TEXT NOT NULL,
  status TEXT NOT NULL,
  result_json TEXT,
  error_json TEXT,
  trace_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteRunStore) CreateRun(ctx context.Context, run RunRecord) error {
	traceJSON, err := json.Marshal(run.Trace)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO runs(id, query, status, trace_json, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)`,
		run.ID,
		run.Query,
		run.Status,
		string(traceJSON),
		run.CreatedAt.Format(time.RFC3339Nano),
		run.UpdatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteRunStore) CompleteRun(ctx context.Context, runID string, status RunStatus, result *RunResult, trace []TraceStep, runErr *ErrorInfo) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	errorJSON, err := json.Marshal(runErr)
	if err != nil {
		return err
	}
	traceJSON, err := json.Marshal(trace)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`UPDATE runs SET status = ?, result_json = ?, error_json = ?, trace_json = ?, updated_at = ? WHERE id = ?`,
		status,
		nullableJSON(resultJSON),
		nullableJSON(errorJSON),
		string(traceJSON),
		time.Now().UTC().Format(time.RFC3339Nano),
		runID,
	)
	return err
}

func (s *SQLiteRunStore) GetRun(ctx context.Context, runID string) (RunRecord, error) {
	var record RunRecord
	var resultJSON sql.NullString
	var errorJSON sql.NullString
	var traceJSON string
	var createdAt string
	var updatedAt string

	err := s.db.QueryRowContext(ctx,
		`SELECT id, query, status, result_json, error_json, trace_json, created_at, updated_at FROM runs WHERE id = ?`, runID).
		Scan(&record.ID, &record.Query, &record.Status, &resultJSON, &errorJSON, &traceJSON, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RunRecord{}, apperrors.NotFound("run not found", map[string]string{"id": runID})
		}
		return RunRecord{}, apperrors.Internal("failed to retrieve run", err)
	}

	if resultJSON.Valid {
		var result RunResult
		if err := json.Unmarshal([]byte(resultJSON.String), &result); err != nil {
			return RunRecord{}, fmt.Errorf("decode result: %w", err)
		}
		record.Result = &result
	}
	if errorJSON.Valid {
		var runErr ErrorInfo
		if err := json.Unmarshal([]byte(errorJSON.String), &runErr); err != nil {
			return RunRecord{}, fmt.Errorf("decode run error: %w", err)
		}
		record.Error = &runErr
	}
	if err := json.Unmarshal([]byte(traceJSON), &record.Trace); err != nil {
		return RunRecord{}, fmt.Errorf("decode trace: %w", err)
	}
	parsedCreatedAt, _ := time.Parse(time.RFC3339Nano, createdAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339Nano, updatedAt)
	record.CreatedAt = parsedCreatedAt
	record.UpdatedAt = parsedUpdatedAt
	return record, nil
}

func (s *SQLiteRunStore) Close() error {
	return s.db.Close()
}

func nullableJSON(value []byte) interface{} {
	if string(value) == "null" {
		return nil
	}
	return string(value)
}
