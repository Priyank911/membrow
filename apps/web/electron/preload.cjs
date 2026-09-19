const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("membrowDesktop", {
  isDesktop: true,
  capturePage: (webContentsId) =>
    ipcRenderer.invoke("desktop:capture-page", webContentsId),
  extractImage: (imageData, apiKey, model) =>
    ipcRenderer.invoke("desktop:extract-image", imageData, apiKey, model),
  mcpConnect: (config) => ipcRenderer.invoke("desktop:mcp-connect", config),
  mcpStore: (item, bucket, config) =>
    ipcRenderer.invoke("desktop:mcp-store", item, bucket, config),
});
