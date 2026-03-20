import { app, BrowserWindow, shell, ipcMain, nativeTheme } from 'electron'
import { join } from 'path'
import Store from 'electron-store'

const store = new Store()

let mainWindow: BrowserWindow | null = null

const isDev = process.env.NODE_ENV === 'development' || !app.isPackaged

function createWindow() {
  mainWindow = new BrowserWindow({
    width:           1280,
    height:          800,
    minWidth:        900,
    minHeight:       600,
    title:           'NetWatch',
    backgroundColor: '#0f172a',
    titleBarStyle:   process.platform === 'darwin' ? 'hiddenInset' : 'default',
    webPreferences: {
      preload:             join(__dirname, 'preload.js'),
      contextIsolation:    true,
      nodeIntegration:     false,
      sandbox:             false,
      webSecurity:         !isDev,
    },
  })

  // Restaurar tamanho/posição da última sessão
  const bounds = store.get('windowBounds') as Electron.Rectangle | undefined
  if (bounds) {
    mainWindow.setBounds(bounds)
  }

  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(join(__dirname, '../dist/index.html'))
  }

  // Abrir links externos no browser do sistema
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url)
    return { action: 'deny' }
  })

  mainWindow.on('close', () => {
    if (mainWindow) {
      store.set('windowBounds', mainWindow.getBounds())
    }
  })

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

app.whenReady().then(() => {
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// ---- IPC Handlers ----

ipcMain.handle('app:get-version', () => app.getVersion())

ipcMain.handle('app:toggle-theme', () => {
  const current = nativeTheme.themeSource
  nativeTheme.themeSource = current === 'dark' ? 'light' : 'dark'
  return nativeTheme.themeSource
})

ipcMain.handle('store:get', (_event, key: string) => store.get(key))
ipcMain.handle('store:set', (_event, key: string, value: unknown) => store.set(key, value))
ipcMain.handle('store:delete', (_event, key: string) => store.delete(key))
