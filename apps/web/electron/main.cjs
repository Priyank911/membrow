const { app, BrowserWindow, ipcMain, session, protocol } = require("electron");
const path = require("path");
const fs = require("fs");

let mainWindow = null;
const mcpSessions = new Map();
const appRoot = path.resolve(__dirname, "..");
const rendererRoot = path.join(appRoot, "build");
const rendererEntry = path.join(rendererRoot, "index.html");

// Keep the custom scheme as a local file-backed origin. This maps SvelteKit's
// root-relative /_app/* imports into the packaged build directory.
protocol.registerSchemesAsPrivileged([
  {
    scheme: "app",
    privileges: {
      standard: true,
      secure: true,
      supportFetchAPI: true,
      corsEnabled: true,
    },
  },
]);

function logStartup(message, error) {
  const detail = error ? ` ${error.stack || error.message || error}` : "";
  console.error(`[startup] ${message}${detail}`);
}

function registerAppProtocol() {
  protocol.registerFileProtocol("app", (request, callback) => {
    try {
      const requestedPath = decodeURIComponent(new URL(request.url).pathname || "/");
      const relativePath = requestedPath.replace(/^[/\\]+/, "") || "index.html";
      const filePath = path.resolve(rendererRoot, relativePath);

      if (filePath !== rendererRoot && !filePath.startsWith(`${rendererRoot}${path.sep}`)) {
        logStartup(`Blocked invalid renderer path: ${requestedPath}`);
        callback({ error: -6 });
        return;
      }

      if (!fs.existsSync(filePath) || fs.statSync(filePath).isDirectory()) {
        // Only extensionless client-side routes use the SPA shell.
        if (path.extname(relativePath) || !fs.existsSync(rendererEntry)) {
          logStartup(`Renderer asset not found: ${filePath}`);
          callback({ error: -6 });
          return;
        }
        callback({ path: rendererEntry });
        return;
      }

      callback({ path: filePath });
    } catch (error) {
      logStartup("Renderer protocol error", error);
      callback({ error: -2 });
    }
  });
}

function createWindow() {
  if (!fs.existsSync(rendererEntry)) {
    throw new Error(`Packaged renderer is missing: ${rendererEntry}`);
  }

  mainWindow = new BrowserWindow({
    width: 1366,
    height: 860,
    minWidth: 900,
    minHeight: 600,
    backgroundColor: "#09090b",
    autoHideMenuBar: true,
    title: "Membrow - Developer Knowledge Browser",
    webPreferences: {
      preload: path.join(__dirname, "preload.cjs"),
      nodeIntegration: false,
      contextIsolation: true,
      webviewTag: true,
      sandbox: false,
    },
  });

  // Handle webview permission requests (camera, mic, etc.) safely
  session.defaultSession.setPermissionRequestHandler(
    (webContents, permission, callback) => {
      callback(false);
    },
  );

  const isDev =
    !app.isPackaged &&
    (process.env.VITE_DEV_SERVER_URL || process.env.NODE_ENV === "development");

  if (isDev) {
    const devServerUrl =
      process.env.VITE_DEV_SERVER_URL || "http://localhost:5173";

    // Load URL
    mainWindow.loadURL(devServerUrl).catch(() => {
      // If dev server is not running yet, retry
      setTimeout(() => {
        mainWindow.loadURL(devServerUrl);
      }, 1500);
    });
  } else {
    // The protocol maps / to build/index.html. Keep the URL at / because
    // SvelteKit's generated route table contains /, not /index.html.
    mainWindow.loadURL("app://local/").catch((error) => {
      logStartup("Could not load packaged renderer", error);
      showStartupError(error);
    });
  }

  mainWindow.webContents.on("did-fail-load", (_event, errorCode, errorDescription, validatedURL) => {
    if (validatedURL.startsWith("app://")) {
      logStartup(`Renderer failed to load (${errorCode}) ${errorDescription} at ${validatedURL}`);
      showStartupError(new Error(`${errorDescription} (${errorCode})`));
    }
  });

  mainWindow.webContents.on("render-process-gone", (_event, details) => {
    logStartup(`Renderer process exited: ${details.reason}`);
  });

  mainWindow.on("closed", () => {
    mainWindow = null;
  });
}

function showStartupError(error) {
  if (!mainWindow || mainWindow.isDestroyed()) return;
  const message = String(error?.message || error || "Unknown startup error").replace(/[<>&\"']/g, (char) => ({
    "<": "&lt;", ">": "&gt;", "&": "&amp;", '"': "&quot;", "'": "&#39;",
  }[char]));
  const html = `<!doctype html><meta charset="utf-8"><body style="background:#09090b;color:#fafafa;font:16px system-ui;padding:40px"><h1>Membrow could not start</h1><p>${message}</p><p>Please restart the app. If this keeps happening, send the startup log to support.</p></body>`;
  mainWindow.loadURL(`data:text/html;charset=utf-8,${encodeURIComponent(html)}`).catch(() => {});
}

// IPC Handlers
ipcMain.handle("desktop:capture-page", async (event, webContentsId) => {
  try {
    const { webContents } = require("electron");
    const targetContents = webContentsId
      ? webContents.fromId(webContentsId)
      : mainWindow.webContents;
    if (!targetContents) {
      throw new Error("Target webview not found");
    }

    // Capture only the visible webview surface. A layout-sized,
    // beyond-viewport capture can stitch repeated content on long pages.
    if (typeof targetContents.capturePage === "function") {
      const image = await targetContents.capturePage();
      return image.toDataURL();
    }

    const debuggerAttached = targetContents.debugger.isAttached();
    if (!debuggerAttached) targetContents.debugger.attach("1.3");
    try {
      await targetContents.debugger.sendCommand("Page.enable");
      const screenshot = await targetContents.debugger.sendCommand(
        "Page.captureScreenshot",
        {
          format: "png",
          captureBeyondViewport: false,
          fromSurface: true,
        },
      );
      return `data:image/png;base64,${screenshot.data}`;
    } finally {
      if (!debuggerAttached) targetContents.debugger.detach();
    }
  } catch (err) {
    console.error("Error in desktop:capture-page:", err);
    throw err;
  }
});

ipcMain.handle("desktop:extract-image", async (event, imageData, apiKey, model) => {
  apiKey = apiKey || process.env.GROQ_API_KEY;
  if (!apiKey) {
    return { success: false, error: "GROQ_API_KEY is not configured" };
  }
  try {
    const response = await fetch("https://api.groq.com/openai/v1/chat/completions", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${apiKey}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: model || "qwen/qwen3.6-27b",
        temperature: 0.1,
        messages: [
          {
            role: "system",
            content:
              "Extract useful knowledge from the screenshot. Return only JSON with title, author, category (research|tool|agent|model|skill), summary, and tags (string array). Do not invent facts.",
          },
          {
            role: "user",
            content: [
              { type: "text", text: "Extract the visible content from this image." },
              { type: "image_url", image_url: { url: imageData } },
            ],
          },
        ],
      }),
    });
    const responseText = await response.text();
    let payload;
    try {
      payload = JSON.parse(responseText);
    } catch {
      payload = parseGroqEventStream(responseText);
    }
    if (!response.ok) {
      return { success: false, error: payload.error?.message || `Groq HTTP ${response.status}` };
    }
    const content = extractGroqContent(payload);
    const data = parseExtractionJson(content);
    return { success: true, data };
  } catch (error) {
    return { success: false, error: error.message || "Groq extraction failed" };
  }
});

function parseGroqEventStream(text) {
  const events = [];
  for (const block of text.split(/\r?\n\r?\n/)) {
    const data = block
      .split(/\r?\n/)
      .filter((line) => line.startsWith("data:"))
      .map((line) => line.slice(5).trim())
      .join("\n");
    if (!data || data === "[DONE]") continue;
    try {
      events.push(JSON.parse(data));
    } catch {
      events.push({ raw: data });
    }
  }
  return { events };
}

function extractGroqContent(payload) {
  if (typeof payload === "string") return payload;
  if (payload?.choices?.[0]?.message?.content) {
    return payload.choices[0].message.content;
  }
  if (payload?.choices?.[0]?.delta?.content) {
    return payload.choices.map((choice) => choice.delta?.content || "").join("");
  }
  if (Array.isArray(payload?.events)) {
    return payload.events
      .map((event) =>
        event.choices?.[0]?.delta?.content ||
        event.choices?.[0]?.message?.content ||
        event.raw ||
        "",
      )
      .join("");
  }
  return "";
}

function parseExtractionJson(content) {
  if (typeof content !== "string") {
    throw new Error("Groq returned no extraction content");
  }
  let withoutEvents = content.trim();
  if (/^event:|\ndata:/m.test(withoutEvents)) {
    const streamedContent = extractGroqContent(parseGroqEventStream(withoutEvents));
    if (streamedContent && streamedContent !== withoutEvents) {
      withoutEvents = streamedContent.trim();
    }
  }
  const unfenced = withoutEvents.replace(/^```(?:json)?\s*/i, "").replace(/\s*```$/, "");
  try {
    return JSON.parse(unfenced);
  } catch {
    const start = unfenced.indexOf("{");
    const end = unfenced.lastIndexOf("}");
    if (start < 0 || end <= start) {
      throw new Error("Groq returned invalid extraction JSON");
    }
    return JSON.parse(unfenced.slice(start, end + 1));
  }
}

ipcMain.handle("desktop:mcp-connect", async (event, config) => {
  if (!config || !config.serverUrl) {
    return { success: false, error: "No MCP server URL provided" };
  }
  const startTime = Date.now();
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const headers = {
      "Content-Type": "application/json",
      Accept: "application/json, text/event-stream",
      ...(config.apiKey
        ? {
            Authorization: config.apiKey.startsWith("Bearer ")
              ? config.apiKey
              : `Bearer ${config.apiKey}`,
          }
        : {}),
    };
    const res = await fetch(config.serverUrl, {
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
    clearTimeout(timeout);
    if (res.ok) {
      const sessionId = res.headers.get("Mcp-Session-Id");
      if (sessionId) mcpSessions.set(config.serverUrl, sessionId);
      const initializedHeaders = { ...headers };
      if (sessionId) {
        initializedHeaders["Mcp-Session-Id"] = sessionId;
        initializedHeaders["MCP-Protocol-Version"] = "2025-06-18";
      }
      await fetch(config.serverUrl, {
        method: "POST",
        headers: initializedHeaders,
        body: JSON.stringify({
          jsonrpc: "2.0",
          method: "notifications/initialized",
          params: {},
        }),
      });
      return { success: true, latencyMs: Date.now() - startTime };
    } else {
      return {
        success: false,
        error: `MCP Server responded with HTTP ${res.status}`,
      };
    }
  } catch (err) {
    return {
      success: false,
      error: `MCP server unreachable at ${config.serverUrl} (offline)`,
    };
  }
});

ipcMain.handle("desktop:mcp-store", async (event, item, bucket, config) => {
  try {
    if (!config || !config.serverUrl) {
      return { success: false, error: "No MCP server configuration provided" };
    }
    const headers = {
      "Content-Type": "application/json",
      Accept: "application/json, text/event-stream",
      "MCP-Protocol-Version": "2025-06-18",
      ...(config.apiKey
        ? {
            Authorization: config.apiKey.startsWith("Bearer ")
              ? config.apiKey
              : `Bearer ${config.apiKey}`,
          }
        : {}),
      ...(mcpSessions.has(config.serverUrl)
        ? { "Mcp-Session-Id": mcpSessions.get(config.serverUrl) }
        : {}),
    };
    const response = await fetch(config.serverUrl, {
      method: "POST",
      headers,
      body: JSON.stringify({
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
      }),
    });
    if (!response.ok) {
      return {
        success: false,
        error: `Memron store_memory failed with HTTP ${response.status}`,
      };
    }
    const result = parseMcpResponse(await response.text());
    if (result.error) {
      return {
        success: false,
        error: result.error.message || "Memron rejected the memory",
      };
    }
    if (result.result?.isError) {
      return {
        success: false,
        error: extractMcpResultText(result) || "Memron memory_store failed",
      };
    }
    if (!result.result) {
      return { success: false, error: "Memron returned no tool result" };
    }
    return { success: true, id: item.id };
  } catch (err) {
    return { success: false, error: err.message };
  }
});

function parseMcpResponse(text) {
  const trimmed = text.trim();
  try {
    return JSON.parse(trimmed);
  } catch {
    const messages = [];
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
    if (!messages.length) {
      throw new Error("Memron returned an empty MCP response");
    }
    return messages[messages.length - 1];
  }
}

function extractMcpResultText(response) {
  return (response.result?.content || [])
    .map((part) => part.text || "")
    .filter(Boolean)
    .join(" ");
}

const gotSingleInstanceLock = app.requestSingleInstanceLock();
if (!gotSingleInstanceLock) {
  app.quit();
} else {
  app.on("second-instance", () => {
    if (!mainWindow) return;
    if (mainWindow.isMinimized()) mainWindow.restore();
    mainWindow.focus();
  });

  app.whenReady().then(() => {
    registerAppProtocol();
    createWindow();

    app.on("activate", () => {
      if (BrowserWindow.getAllWindows().length === 0) {
        createWindow();
      }
    });
  }).catch((error) => {
    logStartup("Electron startup failed", error);
    app.quit();
  });
}

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
