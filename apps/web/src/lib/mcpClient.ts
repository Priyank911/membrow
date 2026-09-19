import type { KnowledgeItem, McpConnectionConfig } from "./types";

const STORAGE_KEY_CONFIG = "membrow_mcp_config";
const STORAGE_KEY_ITEMS_PREFIX = "membrow_bucket_items_";

const DEFAULT_CONFIG: McpConnectionConfig = {
  serverUrl: "http://localhost:3001/mcp",
  bucketName: "developer-research",
  status: "disconnected",
  groqModel: "qwen/qwen3.6-27b",
  lastPingMs: undefined,
  lastSyncedAt: undefined,
};

export const INITIAL_SEEDS: KnowledgeItem[] = [
  {
    id: "seed-1",
    title: "Attention Is All You Need (Transformer Architecture)",
    url: "https://arxiv.org/abs/1706.03762",
    domain: "arxiv.org",
    author: "Vaswani et al.",
    category: "research",
    summary:
      "The foundational paper introducing the Transformer architecture based entirely on self-attention mechanisms without recurrence or convolutions.",
    tags: ["research", "transformer", "attention", "architecture"],
    createdAt: new Date(Date.now() - 86400000 * 2).toISOString(),
    mcpSynced: true,
    notes: "Core foundation reference for LLMs.",
  },
  {
    id: "seed-2",
    title: "Qwen 2.5 72B Instruct - Open Weights Model",
    url: "https://huggingface.co/Qwen/Qwen2.5-72B-Instruct",
    domain: "huggingface.co",
    author: "Qwen Team",
    category: "model",
    summary:
      "State-of-the-art open-weights language model with strong code, math, and multi-turn instruction following capabilities.",
    tags: ["model", "llm", "open-weights", "hf"],
    createdAt: new Date(Date.now() - 86400000).toISOString(),
    mcpSynced: true,
    notes: "Exceptional reasoning benchmarks.",
  },
  {
    id: "seed-3",
    title: "Model Context Protocol (MCP) TypeScript SDK",
    url: "https://github.com/modelcontextprotocol/typescript-sdk",
    domain: "github.com",
    author: "modelcontextprotocol",
    category: "tool",
    summary:
      "Open protocol standard connecting AI assistants with local and remote resources, tools, and memory buckets.",
    tags: ["tool", "mcp", "protocol", "sdk"],
    createdAt: new Date(Date.now() - 3600000 * 5).toISOString(),
    mcpSynced: true,
    notes: "Used for connecting agent memory buckets.",
  },
  {
    id: "seed-4",
    title: "Autonomous Browser Agent Runtime",
    url: "https://x.com/karpathy/status/1750000000000000000",
    domain: "x.com",
    author: "@karpathy",
    category: "agent",
    summary:
      "Discussion on agentic operating systems: browser automation, persistent memory loops, and verification protocols.",
    tags: ["agent", "automation", "memory-layer"],
    createdAt: new Date(Date.now() - 3600000 * 2).toISOString(),
    mcpSynced: true,
    notes: "Key architectural insight for desktop browser agents.",
  },
  {
    id: "seed-5",
    title: "Antigravity Dynamic Skill Specification",
    url: "https://github.com/google/antigravity/tree/main/skills",
    domain: "github.com",
    author: "google-deepmind",
    category: "skill",
    summary:
      "Self-contained executable instructions, YAML frontmatter, and contextual tooling for complex code refactoring and memory synthesis.",
    tags: ["skill", "antigravity", "agent-skills", "prompts"],
    createdAt: new Date(Date.now() - 1800000).toISOString(),
    mcpSynced: true,
    notes: "Standard specification for agent capabilities.",
  },
];

export class McpMemoryClient {
  private config: McpConnectionConfig;
  private sessionId?: string;

  constructor() {
    this.config = this.loadConfig();
  }

  public getConfig(): McpConnectionConfig {
    return { ...this.config };
  }

  public saveConfig(config: Partial<McpConnectionConfig>): McpConnectionConfig {
    this.config = { ...this.config, ...config };
    if (typeof window !== "undefined" && window.localStorage) {
      window.localStorage.setItem(
        STORAGE_KEY_CONFIG,
        JSON.stringify(this.config),
      );
    }
    return this.config;
  }

  private loadConfig(): McpConnectionConfig {
    if (typeof window !== "undefined" && window.localStorage) {
      const stored = window.localStorage.getItem(STORAGE_KEY_CONFIG);
      if (stored) {
        try {
          const parsed = JSON.parse(stored);
          // Fresh sessions always start in disconnected state until verified
          return {
            ...DEFAULT_CONFIG,
            ...parsed,
            status: "disconnected",
            lastPingMs: undefined,
            errorMessage:
              "Server not connected. Clippings are saved in your Local Bucket.",
          };
        } catch {
          // fallback
        }
      }
    }
    return { ...DEFAULT_CONFIG, status: "disconnected" };
  }

  public getItems(bucketName?: string): KnowledgeItem[] {
    const bucket = bucketName || this.config.bucketName || "developer-research";
    if (typeof window === "undefined" || !window.localStorage) {
      return INITIAL_SEEDS;
    }
    const key = `${STORAGE_KEY_ITEMS_PREFIX}${bucket}`;
    const stored = window.localStorage.getItem(key);
    if (!stored) {
      // Initialize with seed items
      window.localStorage.setItem(key, JSON.stringify(INITIAL_SEEDS));
      return INITIAL_SEEDS;
    }
    try {
      return JSON.parse(stored);
    } catch {
      return [];
    }
  }

  public saveItems(items: KnowledgeItem[], bucketName?: string) {
    const bucket = bucketName || this.config.bucketName || "developer-research";
    if (typeof window !== "undefined" && window.localStorage) {
      const key = `${STORAGE_KEY_ITEMS_PREFIX}${bucket}`;
      window.localStorage.setItem(key, JSON.stringify(items));
    }
  }

  public async addItem(item: KnowledgeItem): Promise<KnowledgeItem> {
    if (this.config.status !== "connected") {
      throw new Error("Memron is not connected. Connect before saving a snapshot.");
    }

    const items = this.getItems(this.config.bucketName);
    const existingIndex = items.findIndex((i) => i.url === item.url);
    if (existingIndex >= 0) {
      items[existingIndex] = { ...item, id: items[existingIndex].id };
    } else {
      items.unshift(item);
    }

    try {
      await this.dispatchToMcp(item);
      item.mcpSynced = true;
      this.saveItems(items, this.config.bucketName);
    } catch (err) {
      console.warn("Failed to sync to remote MCP server:", err);
      throw err;
    }

    return item;
  }

  public async extractImage(
    imageData: string,
  ): Promise<Partial<KnowledgeItem>> {
    if (typeof window === "undefined" || !window.membrowDesktop?.extractImage) {
      throw new Error("Image extraction is available in the desktop app");
    }
    const result = await window.membrowDesktop.extractImage(
      imageData,
      this.config.groqApiKey,
      this.config.groqModel,
    );
    if (!result.success || !result.data) {
      throw new Error(result.error || "Groq image extraction failed");
    }
    return result.data;
  }

  public deleteItem(id: string): void {
    const items = this.getItems(this.config.bucketName).filter(
      (i) => i.id !== id,
    );
    this.saveItems(items, this.config.bucketName);
  }

  public async connect(
    overrideConfig?: Partial<McpConnectionConfig>,
  ): Promise<{ success: boolean; latencyMs?: number; error?: string }> {
    if (overrideConfig) {
      this.saveConfig(overrideConfig);
    }

    this.saveConfig({ status: "connecting", errorMessage: undefined });
    const startTime = performance.now();

    // Check if running in Electron desktop bridge
    if (typeof window !== "undefined" && window.membrowDesktop?.mcpConnect) {
      try {
        const result = await window.membrowDesktop.mcpConnect(this.config);
        if (result.success) {
          const latency =
            result.latencyMs || Math.round(performance.now() - startTime);
          this.saveConfig({
            status: "connected",
            lastPingMs: latency,
            lastSyncedAt: new Date().toISOString(),
            errorMessage: undefined,
          });
          return { success: true, latencyMs: latency };
        } else {
          this.saveConfig({
            status: "disconnected",
            lastPingMs: undefined,
            errorMessage:
              result.error ||
              `MCP Server unreachable at ${this.config.serverUrl}. Using Local Bucket.`,
          });
          return { success: false, error: result.error };
        }
      } catch (err: any) {
        this.saveConfig({
          status: "disconnected",
          lastPingMs: undefined,
          errorMessage:
            err.message ||
            `Cannot connect to ${this.config.serverUrl}. Using Local Bucket.`,
        });
        return { success: false, error: err.message };
      }
    }

    // Direct HTTP/JSON-RPC MCP connection
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 2500);

      const headers: Record<string, string> = {
        "Content-Type": "application/json",
        Accept: "application/json, text/event-stream",
      };
      if (this.config.apiKey) {
        headers.Authorization = this.config.apiKey.startsWith("Bearer ")
          ? this.config.apiKey
          : `Bearer ${this.config.apiKey}`;
      }

      const res = await fetch(`${this.config.serverUrl}`, {
        method: "POST",
        headers,
        body: JSON.stringify({
          jsonrpc: "2.0",
          id: "initialize-1",
          method: "initialize",
          params: {
            protocolVersion: "2025-06-18",
            capabilities: {},
            clientInfo: { name: "membrow", version: "0.1.0" },
          },
        }),
        signal: controller.signal,
      });
      clearTimeout(timeoutId);

      const latencyMs = Math.round(performance.now() - startTime);

      if (!res.ok) {
        throw new Error(`MCP Server HTTP ${res.status}: ${res.statusText}`);
      }

      this.sessionId = res.headers.get("Mcp-Session-Id") || undefined;
      const initializedHeaders = { ...headers };
      if (this.sessionId) {
        initializedHeaders["Mcp-Session-Id"] = this.sessionId;
        initializedHeaders["MCP-Protocol-Version"] = "2025-06-18";
      }
      await fetch(`${this.config.serverUrl}`, {
        method: "POST",
        headers: initializedHeaders,
        body: JSON.stringify({
          jsonrpc: "2.0",
          method: "notifications/initialized",
          params: {},
        }),
      });

      this.saveConfig({
        status: "connected",
        lastPingMs: latencyMs,
        lastSyncedAt: new Date().toISOString(),
        errorMessage: undefined,
      });
      return { success: true, latencyMs };
    } catch (err: any) {
      this.saveConfig({
        status: "disconnected",
        lastPingMs: undefined,
        errorMessage: `Server offline at ${this.config.serverUrl}. Knowledge stored in Local Bucket.`,
      });
      return { success: false, error: err.message || "Server offline" };
    }
  }

  public disconnect(): void {
    this.sessionId = undefined;
    this.saveConfig({
      status: "disconnected",
      lastPingMs: undefined,
      errorMessage: undefined,
    });
  }

  private async dispatchToMcp(item: KnowledgeItem): Promise<boolean> {
    if (typeof window !== "undefined" && window.membrowDesktop?.mcpStore) {
      const res = await window.membrowDesktop.mcpStore(
        item,
        this.config.bucketName,
        this.config,
      );
      if (!res.success) {
        throw new Error(res.error || "Memron rejected the memory");
      }
      return true;
    }

    const payload = {
      jsonrpc: "2.0",
      id: `store-${item.id}`,
      method: "tools/call",
      params: {
        name: "memory_store",
        arguments: {
          bucket: "knowledge",
          title: item.title,
          content: [
            item.summary,
            `Source URL: ${item.url}`,
            `Author: ${item.author}`,
            `Category: ${item.category}`,
            `Domain: ${item.domain}`,
            item.notes ? `Notes: ${item.notes}` : "",
            item.imageSnapshot ? "A full-page snapshot was used for extraction." : "",
          ].filter(Boolean).join("\n"),
          tags: item.tags,
        },
      },
    };

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      Accept: "application/json, text/event-stream",
      "MCP-Protocol-Version": "2025-06-18",
    };
    if (this.sessionId) headers["Mcp-Session-Id"] = this.sessionId;
    if (this.config.apiKey) {
      headers.Authorization = this.config.apiKey.startsWith("Bearer ")
        ? this.config.apiKey
        : `Bearer ${this.config.apiKey}`;
    }

    const res = await fetch(this.config.serverUrl, {
      method: "POST",
      headers,
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      throw new Error(`Memron store_memory failed with HTTP ${res.status}`);
    }
    const response = (await parseMcpResponse(await res.text())) as {
      error?: { message?: string };
    };
    if (response.error) {
      throw new Error(response.error.message || "Memron rejected the memory");
    }
    if ((response as any).result?.isError) {
      throw new Error(extractMcpResultText(response) || "Memron memory_store failed");
    }
    if (!(response as any).result) {
      throw new Error("Memron returned no tool result");
    }
    return true;
  }
}

function extractMcpResultText(response: any): string {
  return (response.result?.content || [])
    .map((part: any) => part.text || "")
    .filter(Boolean)
    .join(" ");
}

async function parseMcpResponse(text: string): Promise<unknown> {
  const trimmed = text.trim();
  try {
    return JSON.parse(trimmed);
  } catch {
    const messages: unknown[] = [];
    for (const block of trimmed.split(/\r?\n\r?\n/)) {
      const data = block
        .split(/\r?\n/)
        .filter((line) => line.startsWith("data:"))
        .map((line) => line.slice(5).trim())
        .join("\n");
      if (!data || data === "[DONE]") continue;
      try {
        messages.push(JSON.parse(data));
      } catch {
        throw new Error(`Memron returned invalid MCP response: ${data.slice(0, 120)}`);
      }
    }
    if (!messages.length) throw new Error("Memron returned an empty MCP response");
    return messages[messages.length - 1];
  }
}

export const mcpClient = new McpMemoryClient();
