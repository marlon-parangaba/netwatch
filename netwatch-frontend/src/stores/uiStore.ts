import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface UIState {
  theme:           'dark' | 'light'
  sidebarCollapsed: boolean
  toggleTheme:     () => void
  toggleSidebar:   () => void
  setSidebar:      (collapsed: boolean) => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set, get) => ({
      theme:           'dark',
      sidebarCollapsed: false,

      toggleTheme: () => {
        const next = get().theme === 'dark' ? 'light' : 'dark'
        document.documentElement.classList.toggle('dark', next === 'dark')
        set({ theme: next })
      },

      toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
      setSidebar:    (collapsed) => set({ sidebarCollapsed: collapsed }),
    }),
    {
      name: 'netwatch-ui',
      onRehydrateStorage: () => (state) => {
        if (state) {
          document.documentElement.classList.toggle('dark', state.theme === 'dark')
        }
      },
    }
  )
)
