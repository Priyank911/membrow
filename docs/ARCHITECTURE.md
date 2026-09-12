# Membrow Architecture

## High-level layout

- `apps/api`: Go/Fiber API, agent orchestration, storage, and providers
- `apps/web`: SvelteKit UI for query input, results, and trace timeline
- `docs`: architecture, workflow, and error handling reference

## System diagram

```
[Browser UI]
    |
    v
[SvelteKit Web App]
    |
    v
[Go Fiber API]
    |
    +--> [Search Provider Interface] ---> [Mock Provider | SerpAPI]
    +--> [Fetcher + Extractor Worker Pool]
    +--> [Synthesizer]
    +--> [SQLite Run Store]
    +--> [Structured Logs + Request IDs]
```

## Scalability direction

- Replace `storage.RunStore` SQLite implementation with Postgres-backed store.
- Split fetch/extract worker pool into queue-backed workers.
- Add Redis-backed short-term caching and run pub/sub for horizontal scaling.
- Keep provider abstractions so search/fetch/extract backends can evolve independently.
