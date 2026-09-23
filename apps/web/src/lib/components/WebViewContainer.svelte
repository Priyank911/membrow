<script lang="ts">
  import { onMount } from 'svelte';
  import html2canvas from 'html2canvas';
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
  $: frameBlockedInWeb = !isElectron && isKnownFrameBlockedUrl(currentUrl);

  function isKnownFrameBlockedUrl(url: string): boolean {
    try {
      const hostname = new URL(url).hostname.toLowerCase();
      return hostname === 'duckduckgo.com' || hostname.endsWith('.duckduckgo.com');
    } catch {
      return false;
    }
  }

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

    // In web mode the home page belongs to this document, so capture its
    // rendered DOM instead of returning a placeholder illustration.
    if (isHomePage) {
      const homeElement = document.querySelector('.home-container') as HTMLElement | null;
      if (homeElement) {
        try {
          const rendered = await html2canvas(homeElement, {
            backgroundColor: '#000000',
            imageTimeout: 15000,
            logging: false,
            scale: Math.min(2, window.devicePixelRatio || 1),
            useCORS: true
          });
          return rendered.toDataURL('image/png');
        } catch (err) {
          console.warn('DOM capture failed, falling back to canvas:', err);
        }
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
      <!-- Simulated Webframe Viewport -->
      <div class="simulated-viewport">
        {#if frameBlockedInWeb || iframeError}
          <div class="frame-fallback">
            <div class="frame-fallback-icon"><ShieldCheck size={22} strokeWidth={1.7} /></div>
            <h2>Open this page directly</h2>
            <p>This site does not allow embedded previews in the web version.</p>
            <a href={currentUrl} target="_blank" rel="noreferrer" class="fallback-link">
              <span>Open in a new tab</span>
              <ArrowUpRight size={13} strokeWidth={2} />
            </a>
          </div>
        {:else}
          <iframe
            src={currentUrl}
            title="Web Content"
            class="preview-iframe"
            sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
            on:error={() => (iframeError = true)}
          ></iframe>
        {/if}
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

  .frame-fallback {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.65rem;
    padding: 2rem;
    color: #a1a1aa;
    text-align: center;
    background: radial-gradient(circle at 50% 42%, #18181b 0, #09090b 38%);
  }

  .frame-fallback-icon {
    display: grid;
    place-items: center;
    width: 44px;
    height: 44px;
    border: 1px solid #3f3f46;
    border-radius: 0.75rem;
    color: #d4d4d8;
    background: #121215;
  }

  .frame-fallback h2 {
    margin: 0.2rem 0 0;
    color: #fafafa;
    font-size: 1rem;
    font-weight: 650;
  }

  .frame-fallback p {
    margin: 0;
    max-width: 22rem;
    font-size: 0.78rem;
    line-height: 1.5;
  }

  .fallback-link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    margin-top: 0.25rem;
    padding: 0.45rem 0.7rem;
    border: 1px solid #3f3f46;
    border-radius: 0.4rem;
    color: #fafafa;
    background: #18181b;
    font-size: 0.72rem;
    font-weight: 600;
    text-decoration: none;
  }

  .fallback-link:hover {
    border-color: #71717a;
    background: #27272a;
  }

</style>
