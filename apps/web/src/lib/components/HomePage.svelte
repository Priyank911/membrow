<script lang="ts">
  import type { KnowledgeCategory, KnowledgeItem } from '../types';
  import {
    Search,
    Database,
    ArrowUpRight,
    ShieldCheck
  } from '@lucide/svelte';

  export let onNavigate: (url: string) => void;
  export let onOpenDrawer: () => void;
  export let knowledgeItems: KnowledgeItem[] = [];
  export let isMcpConnected: boolean = false;
  export let bucketName: string = 'developer-research';

  let searchQuery = '';

  function handleSearchSubmit() {
    const trimmed = searchQuery.trim();
    if (!trimmed) return;
    if (/^https?:\/\//i.test(trimmed)) {
      onNavigate(trimmed);
    } else if (/^([a-z0-9-]+\.)+[a-z]{2,}(\/.*)?$/i.test(trimmed)) {
      onNavigate(`https://${trimmed}`);
    } else {
      onNavigate(`https://duckduckgo.com/?q=${encodeURIComponent(trimmed)}`);
    }
  }

  function getCategoryCount(cat: KnowledgeCategory): number {
    return knowledgeItems.filter((i) => i.category === cat).length;
  }
</script>

<div class="home-container">
  <div class="home-wrapper">
    <!-- Hero Brand Section -->
    <header class="home-hero">
      <div class="hero-identity">
        <div class="brand-logo-wrap">
          <img class="membrow-logo" src="/membrow-logo-transparent.png" alt="Membrow robot globe logo" />
        </div>

        <div class="brand-titles">
          <div class="badge-row">
            <span class="hero-badge">Desktop Knowledge Browser</span>
            <span class="hero-badge mono">MCP v1.0</span>
          </div>
          <h1 class="brand-name">Membrow</h1>
          <p class="hero-kicker">A sharper way to gather what matters.</p>
        </div>
      </div>

      <div class="aim-card">
        <div class="aim-title-bar">
          <ShieldCheck size={14} strokeWidth={1.8} class="aim-icon" />
          <span class="aim-label">01 / The mission</span>
        </div>
        <p class="aim-text">
          Browse the open web, clip the signal, and turn scattered developer intelligence into a persistent <strong>Memron MCP Memory Bucket</strong>.
        </p>
        <div class="aim-footer"><span>RESEARCH · TOOLS · AGENTS · MODELS · SKILLS</span></div>
      </div>
    </header>

    <!-- Prominent Omnibox Search Field -->
    <section class="home-search-section">
      <form class="home-search-box" on:submit|preventDefault={handleSearchSubmit}>
        <span class="search-lens-icon">
          <Search size={18} strokeWidth={2} />
        </span>
        <input
          type="text"
          class="home-search-input"
          bind:value={searchQuery}
          placeholder="Search the web or enter a URL (e.g. x.com, github.com, arxiv.org)..."
        />
        <button type="submit" class="home-search-btn" disabled={!searchQuery.trim()}>
          <span>Search</span>
        </button>
      </form>
    </section>

    <!-- Memory Bucket Quick Summary -->
    <section class="bucket-status-banner">
      <div class="banner-left">
        <div class="bucket-icon-box">
          <Database size={15} strokeWidth={2} />
        </div>
        <div class="bucket-info">
          <div class="bucket-name-row">
            <span class="bucket-title">Bucket: {bucketName}</span>
            <span class="connection-pill {isMcpConnected ? 'connected' : 'offline'}">
              {isMcpConnected ? 'Memron MCP Online' : 'Local Storage Active'}
            </span>
          </div>
          <div class="category-breakdown">
            <span class="cat-stat"><strong>{getCategoryCount('research')}</strong> Research</span>
            <span class="stat-sep">•</span>
            <span class="cat-stat"><strong>{getCategoryCount('tool')}</strong> Tools</span>
            <span class="stat-sep">•</span>
            <span class="cat-stat"><strong>{getCategoryCount('agent')}</strong> Agents</span>
            <span class="stat-sep">•</span>
            <span class="cat-stat"><strong>{getCategoryCount('model')}</strong> Models</span>
            <span class="stat-sep">•</span>
            <span class="cat-stat"><strong>{getCategoryCount('skill')}</strong> Skills</span>
          </div>
        </div>
      </div>

      <button class="open-bucket-btn" on:click={onOpenDrawer}>
        <span>Open Bucket</span>
        <ArrowUpRight size={13} strokeWidth={2} />
      </button>
    </section>

  </div>
</div>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Pixelify+Sans:wght@400;500;600;700&display=swap');

  .home-container {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    background: #000000;
    color: #fafafa;
    padding: clamp(1.25rem, 4vw, 3.5rem) clamp(1rem, 4vw, 4rem);
    display: flex;
    justify-content: center;
    user-select: none;
    font-family: 'Pixelify Sans', 'Segoe UI', sans-serif;
  }

  .home-wrapper {
    width: 100%;
    max-width: 1180px;
    display: flex;
    flex-direction: column;
    gap: 1.35rem;
  }

  /* Hero Section */
  .home-hero {
    display: grid;
    grid-template-columns: minmax(300px, 0.82fr) minmax(0, 1.18fr);
    align-items: stretch;
    gap: 0.85rem;
    position: relative;
    padding: 0.5rem 0 1.35rem;
  }

  .home-hero::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 4.5rem;
    height: 2px;
    background: #fafafa;
  }

  .hero-identity {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .brand-logo-wrap {
    display: flex;
    align-items: flex-start;
    justify-content: flex-start;
  }

  .membrow-logo {
    width: 64px;
    height: 64px;
    object-fit: contain;
    border-radius: 0;
    filter: none;
  }

  .brand-titles {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.3rem;
  }

  .badge-row {
    display: flex;
    align-items: center;
    gap: 0.45rem;
  }

  .hero-badge {
    font-size: 0.6rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #a1a1aa;
    background: #0d0d0f;
    border: 1px solid #35353b;
    padding: 0.18rem 0.48rem;
    border-radius: 9999px;
    white-space: nowrap;
  }

  .brand-name {
    margin: 0;
    font-size: clamp(2.6rem, 5vw, 3.75rem);
    font-weight: 700;
    line-height: 0.9;
    letter-spacing: -0.075em;
    color: #fafafa;
  }

  .hero-kicker {
    margin: 0.15rem 0 0;
    color: #71717a;
    font-size: 0.72rem;
    letter-spacing: 0.01em;
  }

  .aim-card {
    background: #0a0a0b;
    border: 1px solid #27272a;
    border-radius: 0.75rem;
    padding: 1rem 1.15rem;
    min-height: 132px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.55rem;
    text-align: left;
    box-shadow: 12px 12px 0 #080808;
  }

  .aim-title-bar {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    color: #a1a1aa;
  }

  .aim-label {
    font-size: 0.62rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #d4d4d8;
  }

  .aim-text {
    margin: 0;
    max-width: 32rem;
    font-size: clamp(0.88rem, 1.3vw, 1rem);
    color: #d4d4d8;
    line-height: 1.35;
  }

  .aim-text strong {
    color: #fafafa;
  }

  .aim-footer {
    border-top: 1px solid #1f1f23;
    padding-top: 0.5rem;
    color: #52525b;
    font: 600 0.52rem/1.2 'Pixelify Sans', monospace;
    letter-spacing: 0.08em;
  }

  /* Search Section */
  .home-search-section {
    width: 100%;
    max-width: 820px;
  }

  .home-search-box {
    display: flex;
    align-items: center;
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.75rem;
    padding: 0.35rem 0.65rem 0.35rem 0.9rem;
    height: 44px;
    transition: all 0.14s ease;
    box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.4);
  }

  .home-search-box:focus-within {
    border-color: #52525b;
    box-shadow: 0 0 0 3px rgba(82, 82, 91, 0.25);
    background: #18181b;
  }

  .search-lens-icon {
    display: flex;
    align-items: center;
    color: #71717a;
    margin-right: 0.7rem;
    flex-shrink: 0;
  }

  .home-search-input {
    flex: 1;
    border: none;
    background: none;
    color: #fafafa;
    font-size: 0.82rem;
    outline: none;
    font-family: inherit;
    letter-spacing: -0.01em;
  }

  .home-search-input::placeholder {
    color: #52525b;
  }

  .home-search-btn {
    background: #fafafa;
    color: #09090b;
    border: none;
    border-radius: 0.45rem;
    padding: 0.4rem 0.8rem;
    font-size: 0.72rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.12s ease;
    flex-shrink: 0;
  }

  .home-search-btn:hover:not(:disabled) {
    background: #e4e4e7;
  }

  .home-search-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* Bucket Status Banner */
  .bucket-status-banner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.65rem;
    padding: 0.6rem 0.8rem;
    max-width: 820px;
  }

  .banner-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .bucket-icon-box {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 0.375rem;
    background: #18181b;
    border: 1px solid #27272a;
    color: #fafafa;
  }

  .bucket-info {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .bucket-name-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .bucket-title {
    font-size: 0.72rem;
    font-weight: 600;
    color: #fafafa;
    font-family: monospace;
  }

  .connection-pill {
    font-size: 0.625rem;
    font-weight: 600;
    padding: 0.15rem 0.45rem;
    border-radius: 0.25rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .connection-pill.connected {
    background: #22c55e18;
    color: #22c55e;
    border: 1px solid #22c55e33;
  }

  .connection-pill.offline {
    background: #27272a;
    color: #a1a1aa;
    border: 1px solid #3f3f46;
  }

  .category-breakdown {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.6rem;
    color: #71717a;
  }

  .stat-sep {
    color: #3f3f46;
  }

  .open-bucket-btn {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0.3rem 0.55rem;
    color: #a1a1aa;
    font-size: 0.68rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.12s ease;
    white-space: nowrap;
  }

  .open-bucket-btn:hover {
    background: #27272a;
    border-color: #3f3f46;
    color: #fafafa;
  }

  .mono {
    font-family: monospace;
  }

  /* Mobile Responsive adjustments */
  @media (max-width: 768px) {
    .home-hero {
      grid-template-columns: 1fr;
    }

    .bucket-status-banner {
      flex-direction: column;
      align-items: flex-start;
    }

    .open-bucket-btn {
      align-self: flex-end;
    }

    .category-breakdown {
      flex-wrap: wrap;
    }
  }

  @media (max-width: 520px) {
    .home-container {
      padding: 1.5rem 1rem;
    }

    .home-hero {
      grid-template-columns: 1fr;
      padding-top: 0.5rem;
    }

    .hero-identity {
      display: flex;
      align-items: center;
      gap: 0.85rem;
    }

    .brand-logo-wrap {
      flex-shrink: 0;
    }

    .membrow-logo {
      width: 52px;
      height: 52px;
    }

    .brand-name {
      font-size: 2.1rem;
    }

    .hero-kicker {
      font-size: 0.7rem;
    }

    .home-search-box {
      height: 42px;
    }

    .home-search-input {
      font-size: 0.8125rem;
    }
  }
</style>
