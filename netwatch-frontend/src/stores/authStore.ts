import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, TokenPair } from '@/types'
import { authService } from '@/services/auth.service'

interface AuthState {
  user: User | null
  tokens: TokenPair | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null

  login:  (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  clearError: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user:            null,
      tokens:          null,
      isAuthenticated: false,
      isLoading:       false,
      error:           null,

      login: async (email, password) => {
        set({ isLoading: true, error: null })
        try {
          const res = await authService.login({ email, password })
          localStorage.setItem('access_token',  res.tokens.access_token)
          localStorage.setItem('refresh_token', res.tokens.refresh_token)
          set({
            user:            res.user,
            tokens:          res.tokens,
            isAuthenticated: true,
            isLoading:       false,
          })
        } catch (err: unknown) {
          const msg = (err as { response?: { data?: { error?: string } } })
            ?.response?.data?.error ?? 'Erro ao fazer login'
          set({ error: msg, isLoading: false })
          throw err
        }
      },

      logout: async () => {
        try { await authService.logout() } catch { /* ignore */ }
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        set({ user: null, tokens: null, isAuthenticated: false })
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'netwatch-auth',
      partialize: (state) => ({
        user:            state.user,
        tokens:          state.tokens,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

// Atalho para verificar roles
export function useIsAdmin()    { return useAuthStore((s) => s.user?.role === 'admin') }
export function useIsOperator() { return useAuthStore((s) => ['admin', 'operator'].includes(s.user?.role ?? '')) }
