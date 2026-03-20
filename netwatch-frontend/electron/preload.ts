import { contextBridge, ipcRenderer } from 'electron'

// API exposta ao processo renderer via window.electronAPI
contextBridge.exposeInMainWorld('electronAPI', {
  getVersion:  () => ipcRenderer.invoke('app:get-version'),
  toggleTheme: () => ipcRenderer.invoke('app:toggle-theme'),
  store: {
    get:    (key: string)                  => ipcRenderer.invoke('store:get', key),
    set:    (key: string, value: unknown)  => ipcRenderer.invoke('store:set', key, value),
    delete: (key: string)                  => ipcRenderer.invoke('store:delete', key),
  },
})
