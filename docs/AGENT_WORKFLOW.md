# Agent Workflow

Membrow follows this workflow for `/agent/run`:

1. **Plan**
   - Validate inputs and initialize run trace metadata.
2. **Search**
   - Use configured provider abstraction (`mock` default, optional `serpapi`).
3. **Fetch**
   - Retrieve source pages with bounded worker pool and retry handling.
4. **Extract**
   - Normalize fetched content into concise text snippets.
5. **Synthesize**
   - Assemble source snippets into a final answer payload.

Each stage updates trace entries (`running` / `success` / `failed`) to support timeline rendering in the web UI.
