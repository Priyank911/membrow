package tools

import (
	"context"
	"errors"
	"time"

	apperrors "github.com/Priyank911/membrow/apps/api/internal/errors"
)

func WithRetry[T any](ctx context.Context, attempts int, backoff time.Duration, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	if attempts < 1 {
		attempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err
		appErr := apperrors.AsAppError(err)
		if !appErr.Retryable || attempt == attempts {
			break
		}

		select {
		case <-ctx.Done():
			return zero, apperrors.Timeout("operation canceled", ctx.Err())
		case <-time.After(backoff * time.Duration(attempt)):
		}
	}
	if lastErr == nil {
		lastErr = errors.New("retry operation failed")
	}
	return zero, lastErr
}
