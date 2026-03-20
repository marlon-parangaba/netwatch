import api from './api'
import type { DashboardStats, APIResponse } from '@/types'

export const dashboardService = {
  async getStats(): Promise<DashboardStats> {
    const res = await api.get<APIResponse<DashboardStats>>('/dashboard/stats')
    return res.data.data!
  },
}
