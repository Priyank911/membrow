package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Priyank911/membrow/apps/api/internal/agent"
	"github.com/Priyank911/membrow/apps/api/internal/config"
	"github.com/Priyank911/membrow/apps/api/internal/extract"
	"github.com/Priyank911/membrow/apps/api/internal/fetch"
	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/storage"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
)

func setupTestServer(t *testing.T) *Server {
	t.Helper()
	tmp := t.TempDir()
	store, err := storage.NewSQLiteRunStore(filepath.Join(tmp, "test.db"))
	if err != nil {
		t.Fatalf("failed to setup store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})
	orch := agent.NewOrchestrator(search.NewMockProvider(), fetch.NewHTTPFetcher("membrow-test"), extract.NewBasicExtractor(), synth.NewBasicSynth(), 1, 0, 5*time.Millisecond)
	svc := agent.NewService(store, orch, 2*time.Second, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	cfg := config.Config{RequestTimeout: 2 * time.Second, DefaultSearchSize: 3}
	return NewServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), search.NewMockProvider(), svc)
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res, err := server.App().Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
}

func TestAgentRunEndpoint(t *testing.T) {
	server := setupTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/agent/run", strings.NewReader(`{"query":"go"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := server.App().Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
}
