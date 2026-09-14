const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("membrowDesktop", {
  isDesktop: true,
  capturePage: (webContentsId) =>
    ipcRenderer.invoke("desktop:capture-page", webContentsId),
  mcpConnect: (config) => ipcRenderer.invoke("desktop:mcp-connect", config),
  mcpStore: (item, bucket) =>
    ipcRenderer.invoke("desktop:mcp-store", item, bucket),
  mcpList: (bucket) => ipcRenderer.invoke("desktop:mcp-list", bucket),
});
