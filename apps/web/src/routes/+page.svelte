<script lang="ts">
  import { onMount } from 'svelte';
  import type {
    BrowserTab,
    KnowledgeCategory,
    KnowledgeItem,
    McpConnectionConfig,
    ToastMessage
  } from '$lib/types';
  import { mcpClient } from '$lib/mcpClient';

  import BrowserChrome from '$lib/components/BrowserChrome.svelte';
  import TabBar from '$lib/components/TabBar.svelte';
  import WebViewContainer from '$lib/components/WebViewContainer.svelte';
  import SidebarDrawer from '$lib/components/SidebarDrawer.svelte';
  import SnapshotModal from '$lib/components/SnapshotModal.svelte';
  import Toast from '$lib/components/Toast.svelte';

  const HOME_URL = 'membrow://home';

  // Tabs state
  let tabs: BrowserTab[] = [
    {
      id: 'tab-1',
      url: HOME_URL,
      title: 'Membrow Home',
      loading: false,
      canGoBack: false,
      canGoForward: false
    }
  ];
  let activeTabId: string = 'tab-1';

  $: activeTab = tabs.find((t) => t.id === activeTabId) || tabs[0];

  // MCP & Knowledge state
  let mcpConfig: McpConnectionConfig = mcpClient.getConfig();
  let knowledgeItems: KnowledgeItem[] = [];

  // Sidebar Drawer state
  let showDrawer: boolean = false;

  // Snapshot Modal state
  let showSnapshotModal: boolean = false;
  let snapshotImage: string = '';
  let snapshotTitle: string = '';
  let snapshotUrl: string = '';
  let snapshotAuthor: string = '';
  let snapshotCategory: KnowledgeCategory = 'research';
  let snapshotSummary: string = '';
  let snapshotTags: string[] = [];
  let snapshotProcessing = false;
  let snapshotExtracted = false;

  // Webview reference
  let webviewContainerRef: WebViewContainer;

  // Toast notifications
  let toasts: ToastMessage[] = [];

  onMount(() => {
    // Load local bucket items and configuration
    knowledgeItems = mcpClient.getItems();
    mcpConfig = mcpClient.getConfig();
  });

  function showToast(text: string, type: 'success' | 'error' | 'info' = 'info') {
    const id = `toast-${Date.now()}-${Math.random().toString(36).substring(2, 6)}`;
    toasts = [...toasts, { id, text, type }];
    setTimeout(() => {
      toasts = toasts.filter((t) => t.id !== id);
    }, 3200);
  }

  // Tab operations
  function handleSelectTab(id: string) {
    activeTabId = id;
  }

  function handleNewTab() {
    const newId = `tab-${Date.now()}`;
    const newTab: BrowserTab = {
      id: newId,
      url: HOME_URL,
      title: 'New Tab',
      loading: false,
      canGoBack: false,
      canGoForward: false
    };
    tabs = [...tabs, newTab];
    activeTabId = newId;
  }

  function handleCloseTab(id: string) {
    if (tabs.length <= 1) return;
    const index = tabs.findIndex((t) => t.id === id);
    tabs = tabs.filter((t) => t.id !== id);
    if (activeTabId === id) {
      const nextTab = tabs[Math.max(0, index - 1)];
      activeTabId = nextTab.id;
    }
  }

  // Navigation operations
  function handleNavigate(newUrl: string) {
    const isHome = newUrl === HOME_URL || !newUrl;
    const resolvedUrl = isHome ? HOME_URL : newUrl;
    const resolvedTitle = isHome ? 'Membrow Home' : newUrl;

    tabs = tabs.map((t) => {
      if (t.id === activeTabId) {
        return {
          ...t,
          url: resolvedUrl,
          title: resolvedTitle,
          loading: !isHome
        };
      }
      return t;
    });
  }

  function handleBack() {
    showToast('Navigation: Back', 'info');
  }

  function handleForward() {
    showToast('Navigation: Forward', 'info');
  }

  function handleReload() {
    if (activeTab.url === HOME_URL) {
      knowledgeItems = mcpClient.getItems();
      showToast('Home page refreshed', 'info');
      return;
    }
    tabs = tabs.map((t) => {
      if (t.id === activeTabId) {
        return { ...t, loading: true };
      }
      return t;
    });
    setTimeout(() => {
      tabs = tabs.map((t) => {
        if (t.id === activeTabId) {
          return { ...t, loading: false };
        }
        return t;
      });
    }, 600);
  }

  function handleHome() {
    handleNavigate(HOME_URL);
  }

  function handlePageLoaded(title: string, url: string) {
    tabs = tabs.map((t) => {
      if (t.id === activeTabId) {
        return {
          ...t,
          title: url === HOME_URL ? 'Membrow Home' : title || t.title,
          url: url || t.url,
          loading: false
        };
      }
      return t;
    });
  }

  function handleLoadingChange(loading: boolean) {
    tabs = tabs.map((t) => {
      if (t.id === activeTabId) {
        return { ...t, loading };
      }
      return t;
    });
  }

  // Snapshot & Clip Workflow
  async function handleTriggerSnapshot() {
    try {
      showToast('Snapshot captured. Click the image to extract it.', 'info');

      // 1. Capture visual image
      let imgData = '';
      if (webviewContainerRef?.captureScreenshot) {
        imgData = await webviewContainerRef.captureScreenshot();
      }

      snapshotImage = imgData;
      snapshotTitle = '';
      snapshotUrl = activeTab.url;
      snapshotAuthor = '';
      snapshotCategory = 'research';
      snapshotSummary = '';
      snapshotTags = [];
      snapshotExtracted = false;

      showSnapshotModal = true;
    } catch (err: any) {
      showToast(err.message || 'Snapshot failed', 'error');
    }
  }

  async function handleExtractSnapshot() {
    if (snapshotExtracted || snapshotProcessing) return;
    snapshotProcessing = true;
    try {
      if (!mcpConfig.groqApiKey) {
        throw new Error('Add a Groq Vision API key in MCP Connection settings first');
      }
      const visionData = await mcpClient.extractImage(snapshotImage);
      snapshotTitle = visionData.title || snapshotTitle;
      snapshotAuthor = visionData.author || snapshotAuthor;
      snapshotCategory = visionData.category || snapshotCategory;
      snapshotSummary = visionData.summary || snapshotSummary;
      snapshotTags = visionData.tags?.length ? visionData.tags : snapshotTags;
      snapshotExtracted = true;
      showToast('Image extraction complete. Review and save to Memron.', 'success');
    } catch (err: any) {
      showToast(err.message || 'Snapshot failed', 'error');
    } finally {
      snapshotProcessing = false;
    }
  }

  async function handleSaveSnapshot(data: Partial<KnowledgeItem>) {
    snapshotProcessing = true;
    try {
      const liveConfig = mcpClient.getConfig();
      if (liveConfig.status !== 'connected') {
        showToast('Reconnecting to Memron...', 'info');
        const connection = await mcpClient.connect();
        mcpConfig = mcpClient.getConfig();
        if (!connection.success) {
          throw new Error(connection.error || 'Connect to Memron before saving a snapshot.');
        }
      } else {
        mcpConfig = liveConfig;
      }

    const newItem: KnowledgeItem = {
      id: `item-${Date.now()}`,
      title: data.title || activeTab.title,
      url: data.url || activeTab.url,
      domain: new URL(data.url || activeTab.url).hostname.replace(/^www\./, ''),
      author: data.author || 'unknown',
      category: data.category || 'research',
      summary: data.summary || '',
      tags: data.tags || [],
      imageSnapshot: data.imageSnapshot || snapshotImage,
      createdAt: new Date().toISOString(),
      mcpSynced: mcpConfig.status === 'connected',
      notes: data.notes || ''
    };

      await mcpClient.addItem(newItem);
      knowledgeItems = mcpClient.getItems();
      showSnapshotModal = false;
      showToast(`Stored in Memron Bucket [${newItem.category}]`, 'success');
    } catch (err: any) {
      showToast(err.message || 'Memron storage failed', 'error');
    } finally {
      snapshotProcessing = false;
    }
  }

  // MCP Sidebar actions
  async function handleConnectMcp() {
    const res = await mcpClient.connect();
    mcpConfig = mcpClient.getConfig();
    if (!res.success) {
      throw new Error(res.error || 'Failed to connect to MCP server');
    }
  }

  function handleDisconnectMcp() {
    mcpClient.disconnect();
    mcpConfig = mcpClient.getConfig();
    showToast('Disconnected from Memron MCP', 'info');
  }

  function handleSaveMcpConfig(cfg: Partial<McpConnectionConfig>) {
    mcpConfig = mcpClient.saveConfig(cfg);
  }

  function handleDeleteKnowledgeItem(id: string) {
    mcpClient.deleteItem(id);
    knowledgeItems = mcpClient.getItems();
    showToast('Item removed from bucket', 'info');
  }
</script>

<div class="browser-app-shell">
  <!-- Tab Strip -->
  <TabBar
    {tabs}
    {activeTabId}
    onSelectTab={handleSelectTab}
    onCloseTab={handleCloseTab}
    onNewTab={handleNewTab}
  />

  <!-- Browser Chrome Top Bar -->
  <BrowserChrome
    url={activeTab?.url || ''}
    loading={activeTab?.loading || false}
    canGoBack={activeTab?.canGoBack || false}
    canGoForward={activeTab?.canGoForward || false}
    isMcpConnected={mcpConfig.status === 'connected'}
    mcpBucketName={mcpConfig.bucketName}
    onToggleDrawer={() => (showDrawer = !showDrawer)}
    onNavigate={handleNavigate}
    onBack={handleBack}
    onForward={handleForward}
    onReload={handleReload}
    onHome={handleHome}
    onSnapshot={handleTriggerSnapshot}
  />

  <!-- Active Viewport Container -->
  <main class="browser-viewport">
    <WebViewContainer
      bind:this={webviewContainerRef}
      currentUrl={activeTab?.url || HOME_URL}
      onNavigate={handleNavigate}
      onOpenDrawer={() => (showDrawer = true)}
      {knowledgeItems}
      isMcpConnected={mcpConfig.status === 'connected'}
      bucketName={mcpConfig.bucketName}
      onPageLoaded={handlePageLoaded}
      onLoadingChange={handleLoadingChange}
    />
  </main>

  <!-- Sidebar Drawer (Hamburger Menu) -->
  <SidebarDrawer
    bind:show={showDrawer}
    {mcpConfig}
    items={knowledgeItems}
    onClose={() => (showDrawer = false)}
    onSaveConfig={handleSaveMcpConfig}
    onConnectMcp={handleConnectMcp}
    onDisconnectMcp={handleDisconnectMcp}
    onDeleteItem={handleDeleteKnowledgeItem}
    onOpenUrlInTab={(url) => {
      handleNavigate(url);
      showDrawer = false;
    }}
    onShowToast={showToast}
  />

  <!-- Snapshot & Extraction Confirmation Modal -->
  <SnapshotModal
    bind:show={showSnapshotModal}
    imageSnapshot={snapshotImage}
    initialTitle={snapshotTitle}
    initialUrl={snapshotUrl}
    initialAuthor={snapshotAuthor}
    initialCategory={snapshotCategory}
    initialSummary={snapshotSummary}
    initialTags={snapshotTags}
    processing={snapshotProcessing}
    extracted={snapshotExtracted}
    onExtract={handleExtractSnapshot}
    bucketName={mcpConfig.bucketName}
    onSave={handleSaveSnapshot}
    onClose={() => (showSnapshotModal = false)}
  />

  <!-- Toast Notifications -->
  <Toast bind:toasts />
</div>

<style>
  .browser-app-shell {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    background: #09090b;
    color: #fafafa;
    overflow: hidden;
    position: relative;
  }

  .browser-viewport {
    flex: 1;
    position: relative;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
</style>
