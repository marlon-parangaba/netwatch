import api from './api'
import type { AlertRule, AlertEvent, CreateAlertRuleRequest, APIResponse, PaginatedResult } from '@/types'

export const alertService = {
  // Rules
  async listRules(params?: { page?: number; limit?: number }): Promise<PaginatedResult<AlertRule>> {
    const res = await api.get<APIResponse<AlertRule[]>>('/alerts/rules', { params })
    return { data: res.data.data ?? [], meta: res.data.meta! }
  },

  async createRule(data: CreateAlertRuleRequest): Promise<AlertRule> {
    const res = await api.post<APIResponse<AlertRule>>('/alerts/rules', data)
    return res.data.data!
  },

  async updateRule(id: string, data: Partial<CreateAlertRuleRequest>): Promise<AlertRule> {
    const res = await api.put<APIResponse<AlertRule>>(`/alerts/rules/${id}`, data)
    return res.data.data!
  },

  async deleteRule(id: string): Promise<void> {
    await api.delete(`/alerts/rules/${id}`)
  },

  // Events
  async listEvents(params?: {
    page?: number
    limit?: number
    status?: string
    severity?: string
    device_id?: string
  }): Promise<PaginatedResult<AlertEvent>> {
    const res = await api.get<APIResponse<AlertEvent[]>>('/alerts/events', { params })
    return { data: res.data.data ?? [], meta: res.data.meta! }
  },

  async acknowledgeEvent(id: string): Promise<AlertEvent> {
    const res = await api.put<APIResponse<AlertEvent>>(`/alerts/events/${id}/ack`)
    return res.data.data!
  },
}
