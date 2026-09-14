const { app, BrowserWindow, ipcMain, session } = require("electron");
const path = require("path");
const fs = require("fs");

let mainWindow = null;

function createWindow() {
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

  const devServerUrl =
    process.env.VITE_DEV_SERVER_URL || "http://localhost:5173";

  // Load URL
  mainWindow.loadURL(devServerUrl).catch(() => {
    // If dev server is not running yet, retry
    setTimeout(() => {
      mainWindow.loadURL(devServerUrl);
    }, 1500);
  });

  mainWindow.on("closed", () => {
    mainWindow = null;
  });
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
    const nativeImage = await targetContents.capturePage();
    return nativeImage.toDataURL();
  } catch (err) {
    console.error("Error in desktop:capture-page:", err);
    throw err;
  }
});

ipcMain.handle("desktop:mcp-connect", async (event, config) => {
  if (!config || !config.serverUrl) {
    return { success: false, error: "No MCP server URL provided" };
  }
  const startTime = Date.now();
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const res = await fetch(config.serverUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(config.apiKey ? { Authorization: `Bearer ${config.apiKey}` } : {}),
      },
      body: JSON.stringify({
        jsonrpc: "2.0",
        id: "ping",
        method: "tools/list",
        params: {},
      }),
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (res.ok) {
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

const getStorageFilePath = (bucketName) => {
  const userData = app.getPath("userData");
  const safeName = (bucketName || "default").replace(/[^a-zA-Z0-9_-]/g, "_");
  return path.join(userData, `membrow_bucket_${safeName}.json`);
};

ipcMain.handle("desktop:mcp-store", async (event, item, bucket) => {
  try {
    const filePath = getStorageFilePath(bucket);
    let items = [];
    if (fs.existsSync(filePath)) {
      items = JSON.parse(fs.readFileSync(filePath, "utf-8"));
    }
    const existingIndex = items.findIndex((i) => i.url === item.url);
    if (existingIndex >= 0) {
      items[existingIndex] = item;
    } else {
      items.unshift(item);
    }
    fs.writeFileSync(filePath, JSON.stringify(items, null, 2), "utf-8");
    return { success: true, id: item.id };
  } catch (err) {
    return { success: false, error: err.message };
  }
});

ipcMain.handle("desktop:mcp-list", async (event, bucket) => {
  try {
    const filePath = getStorageFilePath(bucket);
    if (fs.existsSync(filePath)) {
      const items = JSON.parse(fs.readFileSync(filePath, "utf-8"));
      return { success: true, items };
    }
    return { success: true, items: [] };
  } catch (err) {
    return { success: false, items: [], error: err.message };
  }
});

app.whenReady().then(() => {
  createWindow();

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
