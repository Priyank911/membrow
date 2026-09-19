# Membrow

Membrow is a lightweight, developer-focused desktop browser application designed for rapid research capture and memory synthesis. It integrates directly with **Memron** via the **Model Context Protocol (MCP)** to clip, extract, and ground developer intelligence (papers, tools, models, agents, and skills) into persistent memory buckets.

## Key Features

- **Actual Desktop Browser**: Runs natively via Electron with an embedded Chromium guest web engine (`<webview>`) to browse any site (X/Twitter, GitHub, ArXiv, Hugging Face, Reddit, Google) without `X-Frame-Options` or CORS restrictions.
- **Top-Corner Snapshot & Clip**: Compact, minimalist one-click button to capture high-res page snapshots and automatically extract structured metadata (Title, Platform, Author, Category, Takeaways, Tags).
- **5 Developer Knowledge Categories**:
  - `Research`: Academic preprints, architectures, ablation studies, and benchmarks.
  - `Tools`: Developer utilities, CLIs, libraries, and open-source packages.
  - `Agents`: Autonomous cognitive frameworks, multi-agent orchestrators, and workflows.
  - `Models`: Open weights, LLMs, vision models, quantizations, and checkpoints.
  - `Skills`: Antigravity skills, function schemas, and system prompt specifications.
- **Hamburger Sidebar Drawer**:
  - **Memron MCP Connection**: Connect to Memron MCP servers via standard JSON-RPC protocol (`tools/call` for `store_memory`, `search_memory`, `list_buckets`) with live connection badges and latency counters.
  - **Knowledge Data Bucket**: Full searchable knowledge feed with category filters, Markdown/JSON export, and tab integration.
- **Anti-AI Minimalist Aesthetic**:
  - Dark slate monochromatic palette with high-contrast precision typography.
  - **Strictly zero emojis** across all interfaces.
  - Crisp SVG vector iconography from `@lucide/svelte`.
- **Omnibox & Tab Strip**: Multi-tab browsing, direct URL navigation, DuckDuckGo search integration, and full history controls (Back, Forward, Reload, Home).

## Architecture

```
membrow/
├─ apps/
│  ├─ web/                       # Desktop Browser Shell (Svelte 5 + Vite + Electron)
│  │  ├─ electron/               # Native Electron main process & IPC preload bridge
│  │  │  ├─ main.cjs             # Window management, webview guest host, native capturePage()
│  │  │  └─ preload.cjs          # Typed membrowDesktop bridge
│  │  └─ src/
│  │     ├─ lib/components/      # BrowserChrome, TabBar, WebViewContainer, SidebarDrawer, SnapshotModal
│  │     ├─ lib/mcpClient.ts     # Memron MCP Client & Persistent Knowledge Bucket
│  │     ├─ lib/extractor.ts     # Auto-classification & metadata extraction pipeline
│  │     └─ routes/+page.svelte  # Unified browser shell
│  └─ api/                       # Go + Fiber backend (optional worker service)
├─ docs/                         # Architecture, workflow, and error handling reference
└─ scripts/
   ├─ desktop.ps1                # One-click desktop launcher (PowerShell)
   └─ dev.ps1                    # Multi-process dev helper
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
