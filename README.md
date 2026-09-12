# Membrow

Membrow is a lightweight, robust, Windows-friendly agentic browser/search scaffold that demonstrates a production-grade workflow: **plan → search → fetch → extract → synthesize**.

## Project overview

This initial scaffold focuses on:

- clear monorepo module boundaries
- resilient API behavior with typed errors and retries
- observable run traces for each workflow stage
- local-first defaults (mock search provider + SQLite) while keeping migration seams for Postgres/Redis

## Architecture

```
membrow/
├─ apps/
│  ├─ api/                   # Go + Fiber backend
│  │  ├─ cmd/api             # API entrypoint
│  │  └─ internal/
│  │     ├─ agent            # orchestration + run service
│  │     ├─ api              # routes + HTTP middleware integration
│  │     ├─ config           # env-driven config
│  │     ├─ errors           # typed errors + mapping
│  │     ├─ observability    # structured logging + request IDs
│  │     ├─ search           # provider abstraction + mock/serpapi
│  │     ├─ fetch            # content retrieval abstraction
│  │     ├─ extract          # normalization pipeline
│  │     ├─ synth            # response assembly
│  │     ├─ tools            # retry/backoff helper
│  │     └─ storage          # SQLite run store interface/impl
│  └─ web/                   # SvelteKit frontend
├─ docs/
│  ├─ ARCHITECTURE.md
│  ├─ ERROR_HANDLING.md
│  └─ AGENT_WORKFLOW.md
├─ docker-compose.yml
└─ scripts/                  # cross-platform dev helpers
```

## API endpoints

- `GET /health`
- `POST /search` (passthrough provider search)
- `POST /agent/run` (starts full workflow asynchronously)
- `GET /agent/run/:id` (returns run status, trace, and result/error)

## Environment variables

### API (`apps/api/.env.example`)

| Variable | Default | Purpose |
|---|---|---|
| `API_PORT` | `8080` | Fiber listen port |
| `LOG_LEVEL` | `info` | Structured log level |
| `REQUEST_TIMEOUT_SECONDS` | `15` | Request-level timeout |
| `AGENT_RUN_TIMEOUT_SECONDS` | `45` | Max duration per agent run |
| `SHUTDOWN_TIMEOUT_SECONDS` | `10` | Graceful shutdown timeout |
| `DEFAULT_SEARCH_LIMIT` | `5` | Default search result count |
| `SQLITE_PATH` | `./data/membrow.db` | SQLite file path |
| `SEARCH_PROVIDER` | `mock` | `mock` or `serpapi` |
| `SERPAPI_KEY` | _(empty)_ | Optional real search provider key |
| `WORKER_POOL_SIZE` | `4` | Fetch/extract concurrency limit |
| `MAX_RETRIES` | `2` | Retry count for retryable calls |
| `RETRY_BACKOFF_MS` | `250` | Base retry backoff |
| `FETCH_USER_AGENT` | `membrow-bot/0.1` | HTTP fetch user-agent |

### Web (`apps/web/.env.example`)

| Variable | Default | Purpose |
|---|---|---|
| `PUBLIC_API_BASE_URL` | `http://localhost:8080` | Base URL for backend API |

## Local setup

### Prerequisites

- Go 1.22+
- Node.js 22+
- npm 10+

### Windows (PowerShell)

```powershell
Copy-Item .\apps\api\.env.example .\apps\api\.env
Copy-Item .\apps\web\.env.example .\apps\web\.env
cd .\apps\api; go mod tidy
cd ..\web; npm install
cd ..\..
.\scripts\dev.ps1
```

### macOS/Linux

```bash
cp ./apps/api/.env.example ./apps/api/.env
cp ./apps/web/.env.example ./apps/web/.env
cd ./apps/api && go mod tidy
cd ../web && npm install
cd ../..
./scripts/dev.sh
```

## Docker development

```bash
docker compose up --build
```

- Web: `http://localhost:5173`
- API: `http://localhost:8080`

## Testing

### API tests

```bash
cd apps/api
go test ./...
```

### Web checks

```bash
cd apps/web
npm run lint
```

## Scaling roadmap

1. **SQLite → Postgres**
   - keep `RunStore` interface, add Postgres-backed implementation
2. **Redis cache + queueing**
   - cache repeated queries and stage intermediate artifacts
   - push fetch/extract work onto queue workers
3. **Distributed workers**
   - move stage execution from in-process goroutines to horizontally scalable workers
4. **Run orchestration hardening**
   - add dead-letter queues, retry budgets, and per-provider circuit breaking
5. **Observability expansion**
   - OpenTelemetry traces and metrics exporters

## Follow-up TODOs

- Add authentication/authorization and user-scoped run history.
- Add richer extract pipeline (readability heuristics + content quality scoring).
- Add streaming run updates (SSE/WebSocket) instead of polling.
- Add Postgres implementation and migration tooling.
- Add integration tests using mocked external providers.
