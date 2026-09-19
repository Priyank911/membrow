<script lang="ts">
  import { onMount } from 'svelte';
  import type { KnowledgeItem } from '../types';
  import {
    Globe,
    ExternalLink,
    ShieldCheck,
    ArrowUpRight
  } from '@lucide/svelte';
  import HomePage from './HomePage.svelte';

  export let currentUrl: string = 'membrow://home';
  export let onNavigate: (url: string) => void;
  export let onOpenDrawer: () => void = () => {};
  export let knowledgeItems: KnowledgeItem[] = [];
  export let isMcpConnected: boolean = false;
  export let bucketName: string = 'developer-research';

  export let onPageLoaded: (title: string, url: string) => void;
  export let onLoadingChange: (loading: boolean) => void;

  let webviewEl: any = null;
  let isElectron = false;
  let iframeError = false;

  $: isHomePage = !currentUrl || currentUrl === 'membrow://home' || currentUrl === 'about:blank';

  onMount(() => {
    isElectron = typeof window !== 'undefined' && !!window.membrowDesktop?.isDesktop;

    if (isElectron && webviewEl) {
      setupWebviewListeners();
    }
  });

  function setupWebviewListeners() {
    if (!webviewEl) return;

    webviewEl.addEventListener('did-start-loading', () => {
      onLoadingChange(true);
    });

    webviewEl.addEventListener('did-stop-loading', () => {
      onLoadingChange(false);
      try {
        const title = webviewEl.getTitle?.() || '';
        const url = webviewEl.getURL?.() || currentUrl;
        onPageLoaded(title, url);
      } catch (err) {
        // ignore
      }
    });

    webviewEl.addEventListener('page-title-updated', (e: any) => {
      if (e.title) {
        onPageLoaded(e.title, webviewEl?.getURL?.() || currentUrl);
      }
    });
  }

  export async function captureScreenshot(): Promise<string> {
    // If running in Electron, use native capturePage
    if (isElectron && window.membrowDesktop?.capturePage) {
      try {
        const webContentsId = webviewEl?.getWebContentsId?.();
        return await window.membrowDesktop.capturePage(webContentsId);
      } catch (err) {
        console.warn('Native capture failed, falling back to canvas:', err);
      }
    }

    // Canvas-based snapshot preview fallback
    const canvas = document.createElement('canvas');
    canvas.width = 800;
    canvas.height = 450;
    const ctx = canvas.getContext('2d');
    if (ctx) {
      // Draw background
      ctx.fillStyle = '#121215';
      ctx.fillRect(0, 0, 800, 450);

      // Draw top simulated bar
      ctx.fillStyle = '#18181b';
      ctx.fillRect(0, 0, 800, 40);

      // Draw text
      ctx.fillStyle = '#a1a1aa';
      ctx.font = '14px monospace';
      ctx.fillText(`Membrow Snapshot Capture | ${new Date().toLocaleTimeString()}`, 20, 25);

      ctx.fillStyle = '#fafafa';
      ctx.font = 'bold 20px -apple-system, BlinkMacSystemFont, sans-serif';
      ctx.fillText(currentUrl, 20, 100);

      ctx.fillStyle = '#71717a';
      ctx.font = '14px -apple-system, BlinkMacSystemFont, sans-serif';
      ctx.fillText('Captured page content metadata & visual viewport frame', 20, 130);

      // Grid decoration
      ctx.strokeStyle = '#27272a';
      ctx.lineWidth = 1;
      for (let x = 0; x < 800; x += 40) {
        ctx.beginPath();
        ctx.moveTo(x, 160);
        ctx.lineTo(x, 450);
        ctx.stroke();
      }
      for (let y = 160; y < 450; y += 40) {
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(800, y);
        ctx.stroke();
      }

      ctx.fillStyle = '#22c55e';
      ctx.fillRect(20, 180, 8, 8);
      ctx.fillStyle = '#fafafa';
      ctx.font = '12px monospace';
      ctx.fillText(isMcpConnected ? 'Verified Memron MCP Node' : 'Membrow Local Bucket Node', 36, 188);
    }
    return canvas.toDataURL('image/png');
  }
</script>

<div class="webview-wrapper">
  {#if isHomePage}
    <!-- Dedicated Membrow Home Page -->
    <HomePage
      {onNavigate}
      {onOpenDrawer}
      {knowledgeItems}
      {isMcpConnected}
      {bucketName}
    />
  {:else if isElectron}
    <!-- Native Electron Webview Tag for external sites -->
    <webview
      bind:this={webviewEl}
      src={currentUrl}
      class="native-webview"
      allowpopups
    ></webview>
  {:else}
    <!-- Web / Dev Mode Fallback for external sites -->
    <div class="web-preview-container">
      <div class="web-preview-header">
        <div class="preview-mode-tag">
          <ShieldCheck size={12} strokeWidth={2} />
          <span>Active Web Session</span>
        </div>
        <div class="preview-url monospace">{currentUrl}</div>
        <a href={currentUrl} target="_blank" rel="noreferrer" class="open-external-link">
          <span>Open Direct</span>
          <ArrowUpRight size={12} strokeWidth={2} />
        </a>
      </div>

      <!-- Simulated Webframe Viewport -->
      <div class="simulated-viewport">
        <iframe
          src={currentUrl}
          title="Web Content"
          class="preview-iframe"
          sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
          on:error={() => (iframeError = true)}
        ></iframe>
      </div>
    </div>
  {/if}
</div>

<style>
  .webview-wrapper {
    flex: 1;
    width: 100%;
    height: 100%;
    position: relative;
    background: #09090b;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .native-webview {
    width: 100%;
    height: 100%;
    border: none;
    background: #09090b;
  }

  .web-preview-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #09090b;
    overflow: hidden;
  }

  .web-preview-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.8rem;
    padding: 0.45rem 1rem;
    background: #121215;
    border-bottom: 1px solid #1f1f23;
  }

  .preview-mode-tag {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.6875rem;
    font-weight: 600;
    color: #22c55e;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    flex-shrink: 0;
  }

  .preview-url {
    font-size: 0.75rem;
    color: #a1a1aa;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 500px;
    flex: 1;
  }

  .open-external-link {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    color: #a1a1aa;
    text-decoration: none;
    font-size: 0.6875rem;
    padding: 0.2rem 0.5rem;
    border-radius: 0.25rem;
    border: 1px solid #27272a;
    transition: all 0.12s ease;
    flex-shrink: 0;
  }

  .open-external-link:hover {
    background: #18181b;
    color: #fafafa;
    border-color: #3f3f46;
  }

  .simulated-viewport {
    flex: 1;
    width: 100%;
    height: 100%;
    background: #09090b;
  }

  .preview-iframe {
    width: 100%;
    height: 100%;
    border: none;
    background: #ffffff;
  }

  .monospace {
    font-family: monospace;
  }

  @media (max-width: 640px) {
    .preview-url {
      max-width: 180px;
    }
  }
</style>
