<script lang="ts">
  import type { RunRecord, SearchItem } from '$lib/types';

  const apiBase = import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

  let query = '';
  let loading = false;
  let passthroughResults: SearchItem[] = [];
  let run: RunRecord | null = null;
  let errorMessage = '';

  async function runSearch() {
    errorMessage = '';
    loading = true;
    run = null;
    try {
      const response = await fetch(`${apiBase}/search`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query, limit: 5 })
      });
      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload?.error?.message || 'Search request failed');
      }
      passthroughResults = payload.data.results;
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unknown error';
    } finally {
      loading = false;
    }
  }

  async function runAgent() {
    errorMessage = '';
    loading = true;
    passthroughResults = [];
    run = null;

    try {
      const startResponse = await fetch(`${apiBase}/agent/run`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query, limit: 5 })
      });
      const startPayload = await startResponse.json();
      if (!startResponse.ok) {
        throw new Error(startPayload?.error?.message || 'Failed to start agent run');
      }

      const runID = startPayload.data.run_id;
      await pollRun(runID);
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unknown error';
    } finally {
      loading = false;
    }
  }

  async function pollRun(runID: string) {
    const maxPolls = 20;
    for (let i = 0; i < maxPolls; i += 1) {
      const response = await fetch(`${apiBase}/agent/run/${runID}`);
      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload?.error?.message || 'Failed to get run status');
      }
      run = payload.data as RunRecord;
      if (run.status === 'succeeded' || run.status === 'failed') {
        return;
      }
      await new Promise((resolve) => setTimeout(resolve, 800));
    }
    throw new Error('Agent run timed out while waiting for completion');
  }
</script>

<main>
  <h1>Membrow</h1>
  <p>Agentic browser/search scaffold (plan → search → fetch → extract → synthesize)</p>

  <div class="search-box">
    <input bind:value={query} placeholder="Search the web" />
    <button on:click={runSearch} disabled={!query || loading}>Search</button>
    <button on:click={runAgent} disabled={!query || loading}>Run Agent</button>
  </div>

  {#if errorMessage}
    <div class="error" role="alert">{errorMessage}</div>
  {/if}

  {#if passthroughResults.length > 0}
    <section>
      <h2>Search Results</h2>
      <ul>
        {#each passthroughResults as item}
          <li>
            <a href={item.url} target="_blank" rel="noreferrer">{item.title}</a>
            <p>{item.snippet}</p>
          </li>
        {/each}
      </ul>
    </section>
  {/if}

  {#if run}
    <section>
      <h2>Agent Timeline ({run.status})</h2>
      <ul class="timeline">
        {#each run.trace as step}
          <li class={step.status}>
            <strong>{step.name}</strong>
            <span>{step.status}</span>
            {#if step.error}<small>{step.error}</small>{/if}
          </li>
        {/each}
      </ul>
    </section>

    {#if run.result?.synthesis}
      <section>
        <h2>Synthesized Answer</h2>
        <pre>{run.result.synthesis.answer}</pre>
      </section>
    {/if}

    {#if run.error}
      <section class="error" role="alert">
        <strong>{run.error.code}</strong>
        <p>{run.error.message}</p>
      </section>
    {/if}
  {/if}
</main>

<style>
  main {
    max-width: 960px;
    margin: 2rem auto;
    font-family: Arial, Helvetica, sans-serif;
    padding: 0 1rem;
  }

  .search-box {
    display: grid;
    grid-template-columns: 1fr auto auto;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  input {
    padding: 0.8rem;
    border: 1px solid #ccc;
    border-radius: 0.4rem;
  }

  button {
    border: none;
    background: #0057d9;
    color: white;
    padding: 0.8rem 1rem;
    border-radius: 0.4rem;
    cursor: pointer;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  .timeline li {
    border-left: 4px solid #ddd;
    margin: 0.4rem 0;
    padding-left: 0.75rem;
    list-style: none;
  }

  .timeline li.success {
    border-left-color: #0c9b37;
  }

  .timeline li.failed {
    border-left-color: #d7263d;
  }

  .error {
    background: #ffe5e8;
    color: #8a1320;
    border: 1px solid #f5b9bf;
    border-radius: 0.4rem;
    padding: 0.75rem;
    margin: 1rem 0;
  }

  pre {
    white-space: pre-wrap;
    background: #f7f7f7;
    padding: 1rem;
    border-radius: 0.4rem;
  }

  @media (max-width: 720px) {
    .search-box {
      grid-template-columns: 1fr;
    }
  }
</style>
