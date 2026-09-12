# Error Handling

## Principles

- All failures are converted to typed `AppError` values.
- HTTP middleware maps `AppError` to consistent status codes and envelopes.
- Error envelope shape:
  - `code`
  - `message`
  - `details`
  - `retryable`
  - `request_id`

## Retry behavior

- Tool and provider calls use retry+backoff helper.
- Retries only happen when `retryable=true`.
- Timeouts and context cancellations are surfaced as `TIMEOUT` errors.

## Run-level failures

- Agent runs persist failures with structured `error` metadata in SQLite.
- Run traces contain per-stage status and error strings for user-facing timelines.
