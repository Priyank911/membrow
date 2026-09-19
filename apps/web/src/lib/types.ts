export type SearchItem = {
  title: string;
  url: string;
  snippet: string;
};

export type TraceStep = {
  name: string;
  status: "pending" | "running" | "success" | "failed";
  started_at?: string;
  ended_at?: string;
  error?: string;
};

export type RunRecord = {
  id: string;
  query: string;
  status: "queued" | "running" | "succeeded" | "failed";
  result?: {
    search_results: SearchItem[];
    synthesis: {
      answer: string;
      sources: Array<{
        title: string;
        url: string;
        snippet: string;
        text: string;
      }>;
    };
  };
  trace: TraceStep[];
  error?: {
    code: string;
    message: string;
    retryable: boolean;
  };
};

export type KnowledgeCategory =
  "research" | "tool" | "agent" | "model" | "skill";

export interface KnowledgeItem {
  id: string;
  title: string;
  url: string;
  domain: string;
  author: string;
  category: KnowledgeCategory;
  summary: string;
  tags: string[];
  imageSnapshot?: string; // base64 or data URL
  createdAt: string;
  mcpSynced: boolean;
  notes?: string;
}

export interface BrowserTab {
  id: string;
  url: string;
  title: string;
  loading: boolean;
  canGoBack: boolean;
  canGoForward: boolean;
  favicon?: string;
}

export type McpConnectionStatus =
  "disconnected" | "connecting" | "connected" | "error";

export interface ToastMessage {
  id: string;
  text: string;
  type?: "success" | "error" | "info";
}

export interface McpConnectionConfig {
  serverUrl: string;
  bucketName: string;
  apiKey?: string;
  groqApiKey?: string;
  groqModel?: string;
  status: McpConnectionStatus;
  lastPingMs?: number;
  lastSyncedAt?: string;
  errorMessage?: string;
}

export interface DesktopBridge {
  isDesktop: boolean;
  capturePage: (webContentsId?: number) => Promise<string>;
  extractImage: (
    imageData: string,
    apiKey?: string,
    model?: string,
  ) => Promise<{ success: boolean; data?: Partial<KnowledgeItem>; error?: string }>;
  mcpConnect: (
    config: McpConnectionConfig,
  ) => Promise<{ success: boolean; latencyMs?: number; error?: string }>;
  mcpStore: (
    item: KnowledgeItem,
    bucket: string,
    config: McpConnectionConfig,
  ) => Promise<{ success: boolean; id: string; error?: string }>;
}

declare global {
  interface Window {
    membrowDesktop?: DesktopBridge;
  }
}
