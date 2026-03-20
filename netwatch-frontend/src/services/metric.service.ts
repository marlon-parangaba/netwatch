import api from './api'
import type { Metric, MetricSummary, MetricHistoryQuery, APIResponse } from '@/types'

export const metricService = {
  async getHistory(deviceId: string, params: MetricHistoryQuery = {}): Promise<Metric[]> {
    const res = await api.get<APIResponse<Metric[]>>(`/metrics/${deviceId}/history`, { params })
    return res.data.data ?? []
  },

  async getSummary(deviceId: string, type: string, from?: string, to?: string): Promise<MetricSummary> {
    const res = await api.get<APIResponse<MetricSummary>>(`/metrics/${deviceId}/summary`, {
      params: { type, from, to },
    })
    return res.data.data!
  },
}
