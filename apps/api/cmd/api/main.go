package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Priyank911/membrow/apps/api/internal/agent"
	"github.com/Priyank911/membrow/apps/api/internal/api"
	"github.com/Priyank911/membrow/apps/api/internal/config"
	"github.com/Priyank911/membrow/apps/api/internal/extract"
	"github.com/Priyank911/membrow/apps/api/internal/fetch"
	"github.com/Priyank911/membrow/apps/api/internal/observability"
	"github.com/Priyank911/membrow/apps/api/internal/search"
	"github.com/Priyank911/membrow/apps/api/internal/storage"
	"github.com/Priyank911/membrow/apps/api/internal/synth"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger(cfg.LogLevel)

	store, err := storage.NewSQLiteRunStore(cfg.DatabasePath)
	if err != nil {
		logger.Error("failed to initialize sqlite store", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	provider := searchProviderFromConfig(cfg, logger)
	workflow := agent.NewOrchestrator(
		provider,
		fetch.NewHTTPFetcher(cfg.FetchUserAgent),
		extract.NewBasicExtractor(),
		synth.NewBasicSynth(),
		cfg.WorkerPoolSize,
		cfg.MaxRetries,
		time.Duration(cfg.RetryBackoffMS)*time.Millisecond,
	)
	agentService := agent.NewService(store, workflow, cfg.AgentRunTimeout, logger)
	server := api.NewServer(cfg, logger, provider, agentService)

	addr := fmt.Sprintf(":%s", cfg.Port)
	go func() {
		logger.Info("api listening", "addr", addr)
		if err := server.Listen(addr); err != nil {
			logger.Error("server stopped", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := agentService.Shutdown(shutdownCtx); err != nil {
		logger.Warn("agent service shutdown warning", "error", err)
	}
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("shutdown complete")
}

func searchProviderFromConfig(cfg config.Config, logger *slog.Logger) search.Provider {
	if cfg.SearchProvider == "serpapi" && cfg.SerpAPIKey != "" {
		logger.Info("using serpapi search provider")
		return search.NewSerpAPIProvider(cfg.SerpAPIKey)
	}
	logger.Info("using mock search provider")
	return search.NewMockProvider()
}
