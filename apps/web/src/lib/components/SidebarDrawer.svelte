<script lang="ts">
  import type { KnowledgeCategory, KnowledgeItem, McpConnectionConfig } from '../types';
  import {
    X,
    Database,
    Cpu,
    Search,
    BookOpen,
    Wrench,
    Bot,
    Boxes,
    Terminal,
    ExternalLink,
    Copy,
    Check,
    Trash2,
    RefreshCw,
    Shield,
    Globe,
    User,
    Tag,
    Layers
  } from '@lucide/svelte';

  export let show: boolean = false;
  export let mcpConfig: McpConnectionConfig;
  export let items: KnowledgeItem[] = [];

  export let onClose: () => void;
  export let onSaveConfig: (cfg: Partial<McpConnectionConfig>) => void;
  export let onConnectMcp: () => Promise<void>;
  export let onDisconnectMcp: () => void;
  export let onDeleteItem: (id: string) => void;
  export let onOpenUrlInTab: (url: string) => void;
  export let onShowToast: (text: string, type?: 'success' | 'error' | 'info') => void;

  let activeTab: 'bucket' | 'mcp' = 'bucket';
  let selectedCategory: 'all' | KnowledgeCategory = 'all';
  let searchQuery = '';

  let serverUrlInput = mcpConfig.serverUrl;
  let bucketNameInput = mcpConfig.bucketName;
  let apiKeyInput = mcpConfig.apiKey || '';
  let groqApiKeyInput = mcpConfig.groqApiKey || '';
  let groqModelInput = mcpConfig.groqModel || 'qwen/qwen3.6-27b';
  let connecting = false;

  $: serverUrlInput = mcpConfig.serverUrl;
  $: bucketNameInput = mcpConfig.bucketName;

  const categories: Array<{ id: 'all' | KnowledgeCategory; label: string; icon: any }> = [
    { id: 'all', label: 'All', icon: Layers },
    { id: 'research', label: 'Research', icon: BookOpen },
    { id: 'tool', label: 'Tools', icon: Wrench },
    { id: 'agent', label: 'Agents', icon: Bot },
    { id: 'model', label: 'Models', icon: Boxes },
    { id: 'skill', label: 'Skills', icon: Terminal }
  ];

  $: filteredItems = items.filter((item) => {
    const matchesCategory = selectedCategory === 'all' || item.category === selectedCategory;
    if (!matchesCategory) return false;

    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (
      item.title.toLowerCase().includes(q) ||
      item.author.toLowerCase().includes(q) ||
      item.summary.toLowerCase().includes(q) ||
      item.domain.toLowerCase().includes(q) ||
      item.tags.some((t) => t.toLowerCase().includes(q))
    );
  });

  async function handleConnect() {
    connecting = true;
    onSaveConfig({
      serverUrl: serverUrlInput.trim(),
      bucketName: bucketNameInput.trim(),
      apiKey: apiKeyInput.trim() || undefined,
      groqApiKey: groqApiKeyInput.trim() || undefined,
      groqModel: groqModelInput.trim() || 'qwen/qwen3.6-27b'
    });
    try {
      await onConnectMcp();
      onShowToast('Connected to Memron MCP Memory Bucket', 'success');
    } catch (err: any) {
      onShowToast(err.message || 'Connection failed', 'error');
    } finally {
      connecting = false;
    }
  }

  function copyMarkdown(item: KnowledgeItem) {
    const md = `### [${item.title}](${item.url})
**Category**: \`${item.category}\` | **Author**: ${item.author} | **Source**: ${item.domain}
**Tags**: ${item.tags.map((t) => `#${t}`).join(' ')}

${item.summary}
`;
    navigator.clipboard.writeText(md);
    onShowToast('Copied Markdown to clipboard', 'success');
  }

  function copyJson(item: KnowledgeItem) {
    navigator.clipboard.writeText(JSON.stringify(item, null, 2));
    onShowToast('Copied JSON to clipboard', 'success');
  }
</script>

{#if show}
  <div
    class="drawer-backdrop"
    on:click={(e) => e.target === e.currentTarget && onClose()}
    on:keydown={(e) => e.key === 'Escape' && onClose()}
    role="presentation"
  >
    <aside class="drawer-panel" aria-label="Memron Memory & Bucket Sidebar">
      <!-- Drawer Header -->
      <div class="drawer-header">
        <div class="header-branding">
          <div class="header-icon">
            <Database size={16} strokeWidth={2} />
          </div>
          <div class="branding-text">
            <h2>Memron Memory Bucket</h2>
            <div class="branding-meta">
              <span class="mcp-badge">MCP Protocol</span>
              <span class="header-status {mcpConfig.status}">
                <span class="status-dot"></span>
                {mcpConfig.status === 'connected' ? 'Live' : 'Offline'}
              </span>
            </div>
          </div>
        </div>
        <button class="close-btn" on:click={onClose} aria-label="Close sidebar">
          <X size={15} strokeWidth={2} />
        </button>
      </div>

      <!-- Segmented Nav -->
      <div class="segmented-nav">
        <button
          class="nav-tab {activeTab === 'bucket' ? 'active' : ''}"
          on:click={() => (activeTab = 'bucket')}
        >
          <Database size={13} strokeWidth={1.8} />
          <span>Knowledge Bucket</span>
          <span class="count-pill">{items.length}</span>
        </button>

        <button
          class="nav-tab {activeTab === 'mcp' ? 'active' : ''}"
          on:click={() => (activeTab = 'mcp')}
        >
          <Cpu size={13} strokeWidth={1.8} />
          <span>MCP Connection</span>
          <span class="status-dot {mcpConfig.status === 'connected' ? 'connected' : ''}"></span>
        </button>
      </div>

      <!-- Drawer Content -->
      <div class="drawer-body">
        {#if activeTab === 'mcp'}
          <!-- MCP CONNECTION SECTION -->
          <div class="mcp-view">
            <!-- Status Card -->
            <div class="status-card {mcpConfig.status}">
              <div class="status-card-header">
                <div class="status-indicator-wrap">
                  <span class="status-pulse {mcpConfig.status}"></span>
                  <span class="status-label">
                    {#if mcpConfig.status === 'connected'}
                      Connected to Memron MCP
                    {:else if mcpConfig.status === 'connecting'}
                      Connecting...
                    {:else if mcpConfig.status === 'error'}
                      Connection Error
                    {:else}
                      Disconnected
                    {/if}
                  </span>
                </div>
                {#if mcpConfig.lastPingMs !== undefined}
                  <span class="ping-badge">{mcpConfig.lastPingMs}ms</span>
                {/if}
              </div>
              <p class="status-card-desc">
                {mcpConfig.errorMessage || 'Memron MCP protocol bridge is ready for persistent knowledge syncing.'}
              </p>
            </div>

            <!-- Configuration Form -->
            <div class="config-form">
              <div class="form-group">
                <label for="mcp-server-url" class="form-label">MCP Server Endpoint</label>
                <input
                  id="mcp-server-url"
                  type="text"
                  class="form-input monospace"
                  bind:value={serverUrlInput}
                  placeholder="http://localhost:3001/mcp"
                />
                <span class="form-hint">Memron MCP HTTP / SSE / Bridge endpoint</span>
              </div>

              <div class="form-group">
                <label for="mcp-bucket-name" class="form-label">Active Memory Bucket</label>
                <input
                  id="mcp-bucket-name"
                  type="text"
                  class="form-input monospace"
                  bind:value={bucketNameInput}
                  placeholder="developer-research"
                />
                <span class="form-hint">Namespace for research, models, tools, and skills</span>
              </div>

              <div class="form-group">
                <label for="mcp-api-key" class="form-label">API Key / Token (Optional)</label>
                <input
                  id="mcp-api-key"
                  type="password"
                  class="form-input"
                  bind:value={apiKeyInput}
                  placeholder="Bearer token or Memron API key"
                />
              </div>

              <div class="form-group">
                <label for="groq-api-key" class="form-label">Groq Vision API Key</label>
                <input
                  id="groq-api-key"
                  type="password"
                  class="form-input"
                  bind:value={groqApiKeyInput}
                  placeholder="gsk_..."
                />
                <span class="form-hint">Used to extract knowledge from the full screenshot</span>
              </div>

              <div class="form-group">
                <label for="groq-model" class="form-label">Groq Extraction Model</label>
                <input
                  id="groq-model"
                  type="text"
                  class="form-input monospace"
                  bind:value={groqModelInput}
                  placeholder="qwen/qwen3.6-27b"
                />
              </div>

              <div class="actions-row">
                <button
                  type="button"
                  class="btn-mcp-primary"
                  on:click={handleConnect}
                  disabled={connecting}
                >
                  <RefreshCw size={13} strokeWidth={2} class={connecting ? 'spin' : ''} />
                  <span>{mcpConfig.status === 'connected' ? 'Re-Sync Connection' : 'Connect to Memron'}</span>
                </button>

                {#if mcpConfig.status === 'connected'}
                  <button
                    type="button"
                    class="btn-mcp-secondary"
                    on:click={onDisconnectMcp}
                  >
                    Disconnect
                  </button>
                {/if}
              </div>
            </div>

            <!-- Exposed MCP Tools -->
            <div class="mcp-tools-box">
              <span class="tools-heading">Exposed MCP Tools</span>
              <div class="tool-item">
                <code>memory_store</code>
                <span>Save research snapshot & metadata</span>
              </div>
              <div class="tool-item">
                <code>search_memory</code>
                <span>Query indexed developer insights</span>
              </div>
              <div class="tool-item">
                <code>list_buckets</code>
                <span>Inspect active knowledge namespaces</span>
              </div>
            </div>
          </div>
        {:else}
          <!-- KNOWLEDGE BUCKET FEED SECTION -->
          <div class="bucket-view">
            <!-- Search & Filter bar -->
            <div class="search-filter-box">
              <div class="search-wrap">
                <Search size={13} strokeWidth={1.8} class="search-icon" />
                <input
                  type="text"
                  class="search-input"
                  bind:value={searchQuery}
                  placeholder="Search research, tools, models, skills..."
                />
                {#if searchQuery}
                  <button class="clear-search" on:click={() => (searchQuery = '')}>
                    <X size={11} strokeWidth={2} />
                  </button>
                {/if}
              </div>

              <!-- Categories -->
              <div class="category-tabs">
                {#each categories as cat}
                  {@const Icon = cat.icon}
                  <button
                    type="button"
                    class="cat-filter-btn {selectedCategory === cat.id ? 'active' : ''}"
                    on:click={() => (selectedCategory = cat.id)}
                  >
                    <svelte:component this={Icon} size={11} strokeWidth={1.8} />
                    <span>{cat.label}</span>
                  </button>
                {/each}
              </div>
            </div>

            <!-- Knowledge Items List -->
            <div class="items-list">
              {#if filteredItems.length === 0}
                <div class="empty-state">
                  <BookOpen size={24} strokeWidth={1.5} class="empty-icon" />
                  <p class="empty-title">No entries found</p>
                  <p class="empty-desc">
                    Use the <strong>Clip</strong> button at the top corner while browsing X, GitHub, or ArXiv to store items.
                  </p>
                </div>
              {:else}
                {#each filteredItems as item (item.id)}
                  <article class="item-card">
                    {#if item.imageSnapshot}
                      <div class="card-thumbnail">
                        <img src={item.imageSnapshot} alt={item.title} />
                      </div>
                    {/if}

                    <div class="card-header">
                      <div class="card-meta">
                        <span class="category-badge {item.category}">
                          {#if item.category === 'research'}
                            <BookOpen size={10} strokeWidth={2} />
                          {:else if item.category === 'tool'}
                            <Wrench size={10} strokeWidth={2} />
                          {:else if item.category === 'agent'}
                            <Bot size={10} strokeWidth={2} />
                          {:else if item.category === 'model'}
                            <Boxes size={10} strokeWidth={2} />
                          {:else}
                            <Terminal size={10} strokeWidth={2} />
                          {/if}
                          <span>{item.category}</span>
                        </span>
                        <span class="domain-text">{item.domain}</span>
                        <span class="author-text">{item.author}</span>
                      </div>

                      <div class="card-actions">
                        <button
                          class="action-btn"
                          on:click={() => copyMarkdown(item)}
                          title="Copy as Markdown"
                        >
                          <Copy size={12} strokeWidth={1.8} />
                        </button>
                        <button
                          class="action-btn"
                          on:click={() => onOpenUrlInTab(item.url)}
                          title="Open in Tab"
                        >
                          <ExternalLink size={12} strokeWidth={1.8} />
                        </button>
                        <button
                          class="action-btn delete"
                          on:click={() => onDeleteItem(item.id)}
                          title="Delete"
                        >
                          <Trash2 size={12} strokeWidth={1.8} />
                        </button>
                      </div>
                    </div>

                    <h3 class="card-title">
                      <a href={item.url} on:click|preventDefault={() => onOpenUrlInTab(item.url)}>
                        {item.title}
                      </a>
                    </h3>

                    <p class="card-summary">{item.summary}</p>

                    {#if item.tags.length > 0}
                      <div class="card-tags">
                        {#each item.tags as tag}
                          <span class="tag-pill">#{tag}</span>
                        {/each}
                      </div>
                    {/if}
                  </article>
                {/each}
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </aside>
  </div>
{/if}

<style>
  .drawer-backdrop {
    position: fixed;
    top: 81px;
    right: 0;
    bottom: 0;
    left: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(3px);
    z-index: 950;
    display: flex;
    justify-content: flex-end;
  }

  .drawer-panel {
    width: 380px;
    max-width: 90vw;
    align-self: stretch;
    margin-left: auto;
    margin-right: 0;
    height: 100%;
    background: #101012;
    border-left: 1px solid #27272a;
    display: flex;
    flex-direction: column;
    box-shadow: -20px 0 40px rgba(0, 0, 0, 0.5);
    animation: slideInFromRight 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .drawer-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.9rem 1rem 0.8rem;
    border-bottom: 1px solid #27272a;
    background: linear-gradient(135deg, #1b1b20 0%, #141417 68%, #111113 100%);
  }

  .header-branding {
    display: flex;
    align-items: center;
    gap: 0.55rem;
  }

  .header-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 0.375rem;
    background: #27272a;
    color: #fafafa;
  }

  .branding-text h2 {
    margin: 0;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #fafafa;
    letter-spacing: -0.01em;
  }

  .mcp-badge {
    font-size: 0.625rem;
    font-weight: 600;
    color: #a1a1aa;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .branding-meta {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.15rem;
  }

  .header-status {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    color: #71717a;
    font-size: 0.58rem;
    font-weight: 600;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }

  .header-status .status-dot {
    width: 5px;
    height: 5px;
  }

  .header-status.connected {
    color: #86efac;
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    border-radius: 0.375rem;
    border: none;
    background: none;
    color: #71717a;
    cursor: pointer;
  }

  .close-btn:hover {
    background: #27272a;
    color: #fafafa;
  }

  .segmented-nav {
    display: flex;
    background: #18181b;
    border-bottom: 1px solid #27272a;
    padding: 0.45rem 0.6rem;
    gap: 0.25rem;
  }

  .nav-tab {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.4rem;
    padding: 0.48rem 0.6rem;
    background: none;
    border: 1px solid transparent;
    border-radius: 0.375rem;
    color: #a1a1aa;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.12s ease;
    position: relative;
  }

  .nav-tab:hover {
    color: #fafafa;
  }

  .nav-tab.active {
    background: #24242a;
    border-color: #45454f;
    color: #fafafa;
  }

  .nav-tab.active::after {
    content: '';
    position: absolute;
    left: 0.7rem;
    right: 0.7rem;
    bottom: -0.45rem;
    height: 2px;
    border-radius: 2px;
    background: #fafafa;
  }

  .count-pill {
    background: #121215;
    padding: 0.1rem 0.35rem;
    border-radius: 0.25rem;
    font-size: 0.6875rem;
    font-family: monospace;
  }

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #71717a;
  }

  .status-dot.connected {
    background: #22c55e;
  }

  .drawer-body {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
  }

  /* MCP View */
  .mcp-view {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .status-card {
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0.8rem 0.9rem;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .status-card.connected {
    border-color: #22c55e44;
  }

  .status-card.error {
    border-color: #ef444444;
  }

  .status-card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .status-indicator-wrap {
    display: flex;
    align-items: center;
    gap: 0.45rem;
  }

  .status-pulse {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #71717a;
  }

  .status-pulse.connected {
    background: #22c55e;
    box-shadow: 0 0 8px #22c55e88;
  }

  .status-pulse.error {
    background: #ef4444;
  }

  .status-label {
    font-size: 0.8125rem;
    font-weight: 600;
    color: #fafafa;
  }

  .ping-badge {
    font-size: 0.6875rem;
    font-family: monospace;
    color: #22c55e;
    background: #22c55e18;
    padding: 0.1rem 0.4rem;
    border-radius: 0.25rem;
  }

  .status-card-desc {
    margin: 0;
    font-size: 0.75rem;
    color: #a1a1aa;
    line-height: 1.4;
  }

  .config-form {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .form-label {
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #a1a1aa;
  }

  .form-input {
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0.45rem 0.65rem;
    color: #fafafa;
    font-size: 0.8125rem;
    font-family: inherit;
    outline: none;
    transition: border-color 0.12s ease;
  }

  .form-input:focus {
    border-color: #52525b;
  }

  .form-hint {
    font-size: 0.6875rem;
    color: #71717a;
  }

  .monospace {
    font-family: monospace;
  }

  .actions-row {
    display: flex;
    gap: 0.5rem;
    margin-top: 0.4rem;
  }

  .btn-mcp-primary {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.4rem;
    padding: 0.5rem 0.8rem;
    border-radius: 0.375rem;
    border: none;
    background: #fafafa;
    color: #09090b;
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .btn-mcp-primary:hover:not(:disabled) {
    background: #e4e4e7;
  }

  .btn-mcp-secondary {
    padding: 0.5rem 0.8rem;
    border-radius: 0.375rem;
    border: 1px solid #27272a;
    background: none;
    color: #ef4444;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
  }

  .mcp-tools-box {
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0.75rem;
    background: #18181b;
    display: flex;
    flex-direction: column;
    gap: 0.45rem;
  }

  .tools-heading {
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #71717a;
  }

  .tool-item {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    font-size: 0.6875rem;
    color: #a1a1aa;
  }

  .tool-item code {
    color: #fafafa;
    font-family: monospace;
    font-size: 0.75rem;
  }

  /* Bucket View */
  .bucket-view {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .search-filter-box {
    padding: 0.75rem 0.85rem;
    border-bottom: 1px solid #27272a;
    background: #18181b;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .search-wrap {
    display: flex;
    align-items: center;
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0 0.5rem;
    height: 30px;
  }

  .search-wrap:focus-within {
    border-color: #52525b;
  }

  .search-icon {
    color: #71717a;
    margin-right: 0.4rem;
  }

  .search-input {
    flex: 1;
    border: none;
    background: none;
    color: #fafafa;
    font-size: 0.75rem;
    outline: none;
  }

  .clear-search {
    background: none;
    border: none;
    color: #71717a;
    cursor: pointer;
    display: flex;
    align-items: center;
    padding: 0;
  }

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .cat-filter-btn {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.25rem 0.45rem;
    border-radius: 0.25rem;
    border: 1px solid #27272a;
    background: #121215;
    color: #a1a1aa;
    font-size: 0.6875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .cat-filter-btn:hover {
    border-color: #3f3f46;
    color: #fafafa;
  }

  .cat-filter-btn.active {
    background: #27272a;
    border-color: #52525b;
    color: #fafafa;
  }

  .items-list {
    flex: 1;
    overflow-y: auto;
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 3rem 1.5rem;
    color: #71717a;
  }

  .empty-icon {
    margin-bottom: 0.5rem;
    color: #52525b;
  }

  .empty-title {
    margin: 0;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #d4d4d8;
  }

  .empty-desc {
    margin: 0.35rem 0 0;
    font-size: 0.75rem;
    color: #71717a;
    line-height: 1.4;
  }

  .item-card {
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.45rem;
    transition: border-color 0.12s ease;
  }

  .item-card:hover {
    border-color: #3f3f46;
  }

  .card-thumbnail {
    width: 100%;
    height: 90px;
    border-radius: 0.375rem;
    overflow: hidden;
    border: 1px solid #27272a;
    background: #09090b;
  }

  .card-thumbnail img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: top;
    display: block;
  }

  .card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .card-meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.35rem;
  }

  .category-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    background: #27272a;
    color: #fafafa;
    border-radius: 0.25rem;
    padding: 0.15rem 0.4rem;
    font-size: 0.625rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.03em;
  }

  .domain-text {
    font-size: 0.6875rem;
    color: #71717a;
    font-family: monospace;
  }

  .author-text {
    font-size: 0.6875rem;
    color: #a1a1aa;
  }

  .card-actions {
    display: flex;
    align-items: center;
    gap: 0.2rem;
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    border-radius: 0.25rem;
    border: none;
    background: none;
    color: #71717a;
    cursor: pointer;
    padding: 0;
  }

  .action-btn:hover {
    background: #27272a;
    color: #fafafa;
  }

  .action-btn.delete:hover {
    color: #ef4444;
  }

  .card-title {
    margin: 0;
    font-size: 0.8125rem;
    font-weight: 600;
    line-height: 1.3;
  }

  .card-title a {
    color: #fafafa;
    text-decoration: none;
  }

  .card-title a:hover {
    text-decoration: underline;
  }

  .card-summary {
    margin: 0;
    font-size: 0.75rem;
    color: #a1a1aa;
    line-height: 1.4;
  }

  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .tag-pill {
    font-size: 0.625rem;
    color: #71717a;
    background: #121215;
    border-radius: 0.25rem;
    padding: 0.1rem 0.35rem;
    font-family: monospace;
  }

  .spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes slideInFromRight {
    from {
      transform: translateX(100%);
    }
    to {
      transform: translateX(0);
    }
  }

  @media (max-width: 520px) {
    .drawer-panel {
      width: 100vw;
      max-width: 100vw;
    }
  }
</style>
