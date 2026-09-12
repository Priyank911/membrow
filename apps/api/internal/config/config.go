package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port              string
	LogLevel          string
	RequestTimeout    time.Duration
	AgentRunTimeout   time.Duration
	ShutdownTimeout   time.Duration
	DefaultSearchSize int
	DatabasePath      string
	SearchProvider    string
	SerpAPIKey        string
	WorkerPoolSize    int
	MaxRetries        int
	RetryBackoffMS    int
	FetchUserAgent    string
}

func Load() Config {
	return Config{
		Port:              getEnv("API_PORT", "8080"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		RequestTimeout:    durationFromSeconds("REQUEST_TIMEOUT_SECONDS", 15),
		AgentRunTimeout:   durationFromSeconds("AGENT_RUN_TIMEOUT_SECONDS", 45),
		ShutdownTimeout:   durationFromSeconds("SHUTDOWN_TIMEOUT_SECONDS", 10),
		DefaultSearchSize: intFromEnv("DEFAULT_SEARCH_LIMIT", 5),
		DatabasePath:      getEnv("SQLITE_PATH", "./data/membrow.db"),
		SearchProvider:    getEnv("SEARCH_PROVIDER", "mock"),
		SerpAPIKey:        getEnv("SERPAPI_KEY", ""),
		WorkerPoolSize:    intFromEnv("WORKER_POOL_SIZE", 4),
		MaxRetries:        intFromEnv("MAX_RETRIES", 2),
		RetryBackoffMS:    intFromEnv("RETRY_BACKOFF_MS", 250),
		FetchUserAgent:    getEnv("FETCH_USER_AGENT", "membrow-bot/0.1"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func intFromEnv(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func durationFromSeconds(key string, fallback int) time.Duration {
	return time.Duration(intFromEnv(key, fallback)) * time.Second
}
