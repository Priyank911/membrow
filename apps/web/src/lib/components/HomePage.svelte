<script lang="ts">
  import type { KnowledgeCategory, KnowledgeItem } from '../types';
  import {
    Search,
    Globe,
    BookOpen,
    Wrench,
    Bot,
    Boxes,
    Terminal,
    Code,
    Database,
    ArrowUpRight,
    Camera,
    ShieldCheck,
    Cpu,
    Sparkles
  } from '@lucide/svelte';

  export let onNavigate: (url: string) => void;
  export let onOpenDrawer: () => void;
  export let knowledgeItems: KnowledgeItem[] = [];
  export let isMcpConnected: boolean = false;
  export let bucketName: string = 'developer-research';

  let searchQuery = '';

  const developerHubs = [
    {
      name: 'X (Twitter)',
      url: 'https://x.com',
      handle: 'x.com',
      desc: 'Real-time AI research releases, developer threads, and agent architectures.',
      category: 'Agent',
      icon: Bot
    },
    {
      name: 'GitHub Trending',
      url: 'https://github.com/trending',
      handle: 'github.com',
      desc: 'Trending open-source agent frameworks, developer CLIs, and MCP servers.',
      category: 'Tool',
      icon: Code
    },
    {
      name: 'ArXiv (cs.AI / cs.LG)',
      url: 'https://arxiv.org/list/cs.AI/recent',
      handle: 'arxiv.org',
      desc: 'Academic preprints, neural architectures, reasoning, and benchmarks.',
      category: 'Research',
      icon: BookOpen
    },
    {
      name: 'Hugging Face',
      url: 'https://huggingface.co/models',
      handle: 'huggingface.co',
      desc: 'Open weights, quantizations, vision models, and checkpoints.',
      category: 'Model',
      icon: Boxes
    },
    {
      name: 'Reddit / r/LocalLLaMA',
      url: 'https://reddit.com/r/LocalLLaMA',
      handle: 'reddit.com',
      desc: 'Community benchmarks, fine-tuning setups, and developer skills.',
      category: 'Skill',
      icon: Terminal
    },
    {
      name: 'DuckDuckGo Search',
      url: 'https://duckduckgo.com',
      handle: 'duckduckgo.com',
      desc: 'Search the open web without trackers.',
      category: 'Search',
      icon: Search
    }
  ];

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
      <!-- Membrow Custom Geometric Logo -->
      <div class="brand-logo-wrap">
        <svg class="membrow-logo" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="4" y="4" width="40" height="40" rx="12" fill="#121215" stroke="#27272a" stroke-width="2" />
          <path d="M14 24L24 14L34 24L24 34L14 24Z" stroke="#fafafa" stroke-width="2" stroke-linejoin="round" />
          <circle cx="24" cy="24" r="4" fill="#fafafa" />
          <path d="M24 14V8M24 40V34M14 24H8M40 24H34" stroke="#52525b" stroke-width="2" stroke-linecap="round" />
        </svg>
      </div>

      <div class="brand-titles">
        <div class="badge-row">
          <span class="hero-badge">Desktop Knowledge Browser</span>
          <span class="hero-badge mono">MCP v1.0</span>
        </div>
        <h1 class="brand-name">Membrow</h1>
      </div>

      <!-- Aim & Mission -->
      <div class="aim-card">
        <div class="aim-title-bar">
          <ShieldCheck size={14} strokeWidth={1.8} class="aim-icon" />
          <span class="aim-label">Membrow Aim & Mission</span>
        </div>
        <p class="aim-text">
          A lightweight developer browser engineered to browse the open web, clip insights, and synthesize intelligence across <strong>Research</strong>, <strong>Tools</strong>, <strong>Agents</strong>, <strong>Models</strong>, and <strong>Skills</strong> into your persistent <strong>Memron MCP Memory Bucket</strong>.
        </p>
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
          placeholder="Search DuckDuckGo or enter URL (e.g. x.com, github.com, arxiv.org)..."
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

    <!-- Developer Intelligence Launchpad -->
    <section class="launchpad-section">
      <div class="section-header">
        <span class="section-title">Developer Intelligence Launchpad</span>
        <span class="section-hint">Click any feed or enter URL above</span>
      </div>

      <div class="hubs-grid">
        {#each developerHubs as hub}
          {@const Icon = hub.icon}
          <div
            class="hub-tile"
            on:click={() => onNavigate(hub.url)}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === 'Enter' && onNavigate(hub.url)}
          >
            <div class="tile-top">
              <div class="tile-icon-box">
                <svelte:component this={Icon} size={16} strokeWidth={1.8} />
              </div>
              <span class="tile-category-tag">{hub.category}</span>
            </div>

            <h3 class="tile-title">{hub.name}</h3>
            <p class="tile-desc">{hub.desc}</p>

            <div class="tile-footer">
              <span class="tile-handle mono">{hub.handle}</span>
              <span class="tile-arrow">
                <ArrowUpRight size={13} strokeWidth={2} />
              </span>
            </div>
          </div>
        {/each}
      </div>
    </section>

    <!-- Workflow Tip -->
    <footer class="home-footer-tip">
      <div class="tip-content">
        <Camera size={14} strokeWidth={1.8} color="#a1a1aa" />
        <span>While browsing any site, click <strong>Clip</strong> in the top-right corner to extract metadata and store in your bucket.</span>
      </div>
    </footer>
  </div>
</div>

<style>
  .home-container {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    background: #09090b;
    color: #fafafa;
    padding: 2.5rem 1.5rem;
    display: flex;
    justify-content: center;
    user-select: none;
  }

  .home-wrapper {
    width: 100%;
    max-width: 820px;
    display: flex;
    flex-direction: column;
    gap: 1.75rem;
  }

  /* Hero Section */
  .home-hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1.1rem;
  }

  .brand-logo-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .membrow-logo {
    width: 54px;
    height: 54px;
    filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.5));
  }

  .brand-titles {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.35rem;
  }

  .badge-row {
    display: flex;
    align-items: center;
    gap: 0.45rem;
  }

  .hero-badge {
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #a1a1aa;
    background: #18181b;
    border: 1px solid #27272a;
    padding: 0.2rem 0.55rem;
    border-radius: 9999px;
  }

  .brand-name {
    margin: 0;
    font-size: 2.25rem;
    font-weight: 700;
    letter-spacing: -0.03em;
    color: #fafafa;
  }

  /* Aim Card */
  .aim-card {
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.75rem;
    padding: 0.9rem 1.25rem;
    max-width: 640px;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    text-align: left;
  }

  .aim-title-bar {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    color: #a1a1aa;
  }

  .aim-label {
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #a1a1aa;
  }

  .aim-text {
    margin: 0;
    font-size: 0.8125rem;
    color: #d4d4d8;
    line-height: 1.5;
  }

  .aim-text strong {
    color: #fafafa;
  }

  /* Search Section */
  .home-search-section {
    width: 100%;
  }

  .home-search-box {
    display: flex;
    align-items: center;
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.75rem;
    padding: 0.35rem 0.65rem 0.35rem 0.9rem;
    height: 48px;
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
    font-size: 0.9375rem;
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
    padding: 0.45rem 1rem;
    font-size: 0.8125rem;
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
    padding: 0.75rem 1rem;
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
    width: 32px;
    height: 32px;
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
    font-size: 0.8125rem;
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
    font-size: 0.6875rem;
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
    padding: 0.35rem 0.65rem;
    color: #a1a1aa;
    font-size: 0.75rem;
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

  /* Launchpad Grid */
  .launchpad-section {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .section-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
  }

  .section-title {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #a1a1aa;
  }

  .section-hint {
    font-size: 0.6875rem;
    color: #71717a;
  }

  .hubs-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.75rem;
  }

  .hub-tile {
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0.85rem;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    cursor: pointer;
    transition: all 0.14s ease;
  }

  .hub-tile:hover {
    background: #18181b;
    border-color: #3f3f46;
    transform: translateY(-1px);
  }

  .tile-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .tile-icon-box {
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

  .tile-category-tag {
    font-size: 0.625rem;
    font-weight: 600;
    text-transform: uppercase;
    color: #a1a1aa;
    background: #18181b;
    padding: 0.15rem 0.4rem;
    border-radius: 0.25rem;
    letter-spacing: 0.03em;
  }

  .tile-title {
    margin: 0;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #fafafa;
  }

  .tile-desc {
    margin: 0;
    font-size: 0.6875rem;
    color: #71717a;
    line-height: 1.4;
    flex: 1;
  }

  .tile-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 0.35rem;
    padding-top: 0.35rem;
    border-top: 1px solid #1c1c21;
  }

  .tile-handle {
    font-size: 0.6875rem;
    color: #52525b;
  }

  .tile-arrow {
    color: #71717a;
    display: flex;
    align-items: center;
  }

  /* Tip footer */
  .home-footer-tip {
    display: flex;
    justify-content: center;
    padding-top: 0.5rem;
  }

  .tip-content {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.75rem;
    color: #71717a;
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 0.45rem 0.85rem;
  }

  .tip-content strong {
    color: #fafafa;
  }

  .mono {
    font-family: monospace;
  }

  /* Mobile Responsive adjustments */
  @media (max-width: 768px) {
    .hubs-grid {
      grid-template-columns: repeat(2, 1fr);
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

    .hubs-grid {
      grid-template-columns: 1fr;
    }

    .brand-name {
      font-size: 1.75rem;
    }

    .home-search-box {
      height: 42px;
    }

    .home-search-input {
      font-size: 0.8125rem;
    }
  }
</style>
