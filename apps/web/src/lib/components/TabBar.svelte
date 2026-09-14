<script lang="ts">
  import type { BrowserTab } from '../types';
  import { Plus, X, Globe, Loader2 } from '@lucide/svelte';

  export let tabs: BrowserTab[] = [];
  export let activeTabId: string = '';

  export let onSelectTab: (id: string) => void;
  export let onCloseTab: (id: string) => void;
  export let onNewTab: () => void;
</script>

<div class="tab-bar">
  <div class="tabs-list">
    {#each tabs as tab (tab.id)}
      <div
        class="tab {tab.id === activeTabId ? 'active' : ''}"
        on:click={() => onSelectTab(tab.id)}
        role="tab"
        aria-selected={tab.id === activeTabId}
        tabindex="0"
        on:keydown={(e) => e.key === 'Enter' && onSelectTab(tab.id)}
      >
        <span class="tab-icon">
          {#if tab.loading}
            <span class="spinner"><Loader2 size={13} /></span>
          {:else}
            <Globe size={13} strokeWidth={1.75} />
          {/if}
        </span>
        <span class="tab-title">{tab.title || 'New Tab'}</span>
        {#if tabs.length > 1}
          <button
            class="tab-close"
            on:click|stopPropagation={() => onCloseTab(tab.id)}
            title="Close tab"
            aria-label="Close tab"
          >
            <X size={12} strokeWidth={2} />
          </button>
        {/if}
      </div>
    {/each}

    <button class="new-tab-btn" on:click={onNewTab} title="New tab" aria-label="New tab">
      <Plus size={14} strokeWidth={2} />
    </button>
  </div>
</div>

<style>
  .tab-bar {
    display: flex;
    align-items: center;
    background: #09090b;
    border-bottom: 1px solid #1f1f23;
    padding: 0.35rem 0.6rem 0;
    overflow-x: auto;
    scrollbar-width: none;
  }

  .tab-bar::-webkit-scrollbar {
    display: none;
  }

  .tabs-list {
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    padding: 0.4rem 0.75rem;
    max-width: 220px;
    min-width: 130px;
    background: #121215;
    color: #a1a1aa;
    border: 1px solid #27272a;
    border-bottom: none;
    border-top-left-radius: 0.5rem;
    border-top-right-radius: 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
    user-select: none;
    transition: background 0.12s ease, color 0.12s ease;
  }

  .tab:hover {
    background: #18181b;
    color: #fafafa;
  }

  .tab.active {
    background: #1c1c21;
    color: #fafafa;
    border-color: #3f3f46;
    border-bottom: 1px solid #1c1c21;
    margin-bottom: -1px;
    z-index: 2;
  }

  .tab-icon {
    display: flex;
    align-items: center;
    flex-shrink: 0;
    color: #71717a;
  }

  .tab.active .tab-icon {
    color: #a1a1aa;
  }

  .spinner {
    display: flex;
    animation: spin 1s linear infinite;
  }

  .tab-title {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    letter-spacing: -0.01em;
  }

  .tab-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    border: none;
    background: none;
    color: #71717a;
    border-radius: 0.25rem;
    cursor: pointer;
    flex-shrink: 0;
    padding: 0;
    transition: background 0.1s ease, color 0.1s ease;
  }

  .tab-close:hover {
    background: #27272a;
    color: #fafafa;
  }

  .new-tab-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    border-radius: 0.375rem;
    border: 1px solid transparent;
    background: none;
    color: #71717a;
    cursor: pointer;
    margin-left: 0.25rem;
    transition: all 0.12s ease;
  }

  .new-tab-btn:hover {
    background: #18181b;
    border-color: #27272a;
    color: #fafafa;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
