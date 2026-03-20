import api from './api'
import type { AuthResponse, LoginRequest, User } from '@/types'

export const authService = {
  async login(data: LoginRequest): Promise<AuthResponse> {
    const res = await api.post<{ success: boolean; data: AuthResponse }>('/auth/login', data)
    return res.data.data!
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout').catch(() => {})
  },

  async refresh(refreshToken: string): Promise<{ access_token: string }> {
    const res = await api.post('/auth/refresh', { refresh_token: refreshToken })
    return res.data.data
  },

  async me(): Promise<User> {
    const res = await api.get<{ success: boolean; data: User }>('/auth/me')
    return res.data.data!
  },
}
