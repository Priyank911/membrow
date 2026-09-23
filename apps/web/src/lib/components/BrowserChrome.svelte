<script lang="ts">
  import {
    Menu,
    ArrowLeft,
    ArrowRight,
    RotateCw,
    Home,
    Search,
    Lock,
    Camera,
    Database,
    ExternalLink
  } from '@lucide/svelte';

  export let url: string = '';
  export let loading: boolean = false;
  export let canGoBack: boolean = false;
  export let canGoForward: boolean = false;
  export let isMcpConnected: boolean = false;
  export let mcpBucketName: string = 'developer-research';

  export let onToggleDrawer: () => void;
  export let onNavigate: (newUrl: string) => void;
  export let onBack: () => void;
  export let onForward: () => void;
  export let onReload: () => void;
  export let onSnapshot: () => void;

  let inputUrl = '';

  $: inputUrl = url === 'membrow://home' ? '' : isSearchUrl(url) ? 'Search results' : url;

  function isSearchUrl(value: string): boolean {
    try {
      const parsed = new URL(value);
      return parsed.hostname === 'duckduckgo.com' || parsed.hostname.endsWith('.duckduckgo.com');
    } catch {
      return false;
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      let trimmed = inputUrl.trim();
      if (!trimmed) {
        onNavigate('membrow://home');
        return;
      }

      // If it looks like a URL or starts with http/https
      if (/^https?:\/\//i.test(trimmed)) {
        onNavigate(trimmed);
      } else if (/^([a-z0-9-]+\.)+[a-z]{2,}(\/.*)?$/i.test(trimmed)) {
        onNavigate(`https://${trimmed}`);
      } else {
        // Search query
        onNavigate(`https://duckduckgo.com/?q=${encodeURIComponent(trimmed)}`);
      }
    }
  }
</script>

<header class="browser-chrome">
  <!-- Left section: Navigation buttons -->
  <div class="nav-controls">
    <div class="history-controls">
      <button
        class="icon-btn"
        disabled={!canGoBack}
        on:click={onBack}
        title="Back"
        aria-label="Back"
      >
        <ArrowLeft size={15} strokeWidth={1.8} />
      </button>

      <button
        class="icon-btn"
        disabled={!canGoForward}
        on:click={onForward}
        title="Forward"
        aria-label="Forward"
      >
        <ArrowRight size={15} strokeWidth={1.8} />
      </button>

      <button
        class="icon-btn"
        on:click={onReload}
        title={loading ? 'Stop' : 'Reload'}
        aria-label="Reload"
      >
        <span class="icon-inner {loading ? 'spin' : ''}">
          <RotateCw size={14} strokeWidth={1.8} />
        </span>
      </button>
    </div>
  </div>

  <!-- Center section: Omnibox (Search & Address Bar) -->
  <div class="omnibox-container">
    <div class="omnibox">
      <span class="omnibox-protocol">
        {#if !url || url === 'membrow://home'}
          <Home size={12} strokeWidth={2} color="#a1a1aa" />
        {:else if url.startsWith('https://')}
          <Lock size={12} strokeWidth={2} color="#22c55e" />
        {:else}
          <Search size={12} strokeWidth={2} color="#71717a" />
        {/if}
      </span>
      <input
        type="text"
        class="omnibox-input"
        bind:value={inputUrl}
        on:keydown={handleKeydown}
        placeholder="Enter a URL or search the web..."
        spellcheck="false"
      />
      {#if loading}
        <div class="loading-line"></div>
      {/if}
    </div>
  </div>

  <!-- Right section: Memory, capture, and menu utilities -->
  <div class="action-controls">
    <button
      class="compact-badge-btn"
      on:click={onToggleDrawer}
      title={isMcpConnected ? `Memron MCP Connected: ${mcpBucketName}` : `Local Storage Bucket: ${mcpBucketName} (MCP Offline)`}
      aria-label={isMcpConnected ? `Memron MCP Connected: ${mcpBucketName}` : `Local Storage Bucket: ${mcpBucketName} (MCP Offline)`}
    >
      <Database size={12} strokeWidth={1.75} />
      <span class="status-indicator {isMcpConnected ? 'online' : 'offline'}"></span>
    </button>

    <button
      class="snapshot-btn"
      on:click={onSnapshot}
      title="Clip snapshot & extract to Memron Memory Bucket"
      aria-label="Snapshot & Clip"
    >
      <Camera size={14} strokeWidth={2} />
    </button>

    <button
      class="icon-btn hamburger-btn"
      on:click={onToggleDrawer}
      title="Menu & Memron MCP Memory"
      aria-label="Toggle Menu"
    >
      <Menu size={16} strokeWidth={1.8} />
      {#if isMcpConnected}
        <span class="mcp-dot-badge" title="Memron MCP Connected"></span>
      {/if}
    </button>
  </div>
</header>

<style>
  .browser-chrome {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.55rem;
    padding: 0.4rem 0.65rem;
    background: #18181b;
    border-bottom: 1px solid #27272a;
    position: relative;
    user-select: none;
  }

  .nav-controls {
    display: flex;
    align-items: center;
    gap: 0.15rem;
    flex-shrink: 0;
  }

  .history-controls {
    display: flex;
    align-items: center;
    gap: 0.15rem;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 0.375rem;
    border: 1px solid transparent;
    background: none;
    color: #a1a1aa;
    cursor: pointer;
    transition: all 0.12s ease;
    padding: 0;
  }

  .icon-btn:hover:not(:disabled) {
    background: #1c1c21;
    border-color: #27272a;
    color: #fafafa;
  }

  .icon-btn:disabled {
    opacity: 0.35;
    cursor: not-allowed;
  }

  .hamburger-btn {
    position: relative;
    color: #d4d4d8;
    border-color: #3f3f46;
    background: #27272a;
  }

  .hamburger-btn:hover {
    background: #3f3f46;
    border-color: #52525b;
  }

  .mcp-dot-badge {
    position: absolute;
    top: 5px;
    right: 5px;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #22c55e;
    box-shadow: 0 0 6px #22c55e88;
  }

  .omnibox-container {
    flex: 1;
    min-width: 0;
    position: relative;
  }

  .omnibox {
    display: flex;
    align-items: center;
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0 0.6rem;
    height: 32px;
    transition: border-color 0.12s ease, box-shadow 0.12s ease;
    position: relative;
    overflow: hidden;
  }

  .omnibox:focus-within {
    border-color: #52525b;
    box-shadow: 0 0 0 2px rgba(82, 82, 91, 0.2);
  }

  .omnibox-protocol {
    display: flex;
    align-items: center;
    color: #71717a;
    margin-right: 0.45rem;
    flex-shrink: 0;
  }

  .icon-inner {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .omnibox-input {
    flex: 1;
    border: none;
    background: none;
    color: #fafafa;
    font-size: 0.8125rem;
    outline: none;
    font-family: inherit;
    letter-spacing: -0.01em;
  }

  .omnibox-input::placeholder {
    color: #52525b;
  }

  .loading-line {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: #fafafa;
    animation: loadingBar 1.2s infinite ease-in-out;
  }

  .action-controls {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    flex-shrink: 0;
  }

  .compact-badge-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 28px;
    height: 28px;
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0;
    color: #a1a1aa;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .compact-badge-btn:hover {
    background: #222227;
    border-color: #3f3f46;
    color: #fafafa;
  }

  .status-indicator {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    position: absolute;
    right: 4px;
    top: 4px;
  }

  .status-indicator.online {
    background: #22c55e;
    box-shadow: 0 0 5px #22c55e66;
  }

  .status-indicator.offline {
    background: #71717a;
  }

  .snapshot-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: #fafafa;
    color: #09090b;
    border: none;
    border-radius: 0.45rem;
    padding: 0;
    cursor: pointer;
    transition: all 0.12s ease;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  }

  .snapshot-btn:hover {
    background: #e4e4e7;
    transform: translateY(-0.5px);
  }

  .snapshot-btn:active {
    transform: translateY(0);
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

  @keyframes loadingBar {
    0% {
      transform: translateX(-100%);
    }
    50% {
      transform: translateX(0%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  @media (max-width: 768px) {
    .omnibox-container {
      max-width: none;
    }
  }

  @media (max-width: 540px) {
    .browser-chrome {
      padding: 0.35rem 0.45rem;
      gap: 0.35rem;
    }

    .history-controls .icon-btn:not(:first-child) {
      display: none; /* Hide forward and reload on very narrow mobile screens */
    }

    .omnibox {
      height: 30px;
      padding: 0 0.45rem;
    }

    .omnibox-input {
      font-size: 0.75rem;
    }
  }
</style>
