import api from './api'
import type {
  Device,
  CreateDeviceRequest,
  UpdateDeviceRequest,
  DiscoverRequest,
  DiscoverResult,
  PaginatedResult,
  APIResponse,
} from '@/types'

interface DeviceListParams {
  page?: number
  limit?: number
  status?: string
  type?: string
  search?: string
  enabled?: boolean
}

export const deviceService = {
  async list(params: DeviceListParams = {}): Promise<PaginatedResult<Device>> {
    const res = await api.get<APIResponse<Device[]>>('/devices', { params })
    return {
      data: res.data.data ?? [],
      meta: res.data.meta!,
    }
  },

  async getById(id: string): Promise<Device> {
    const res = await api.get<APIResponse<Device>>(`/devices/${id}`)
    return res.data.data!
  },

  async create(data: CreateDeviceRequest): Promise<Device> {
    const res = await api.post<APIResponse<Device>>('/devices', data)
    return res.data.data!
  },

  async update(id: string, data: UpdateDeviceRequest): Promise<Device> {
    const res = await api.put<APIResponse<Device>>(`/devices/${id}`, data)
    return res.data.data!
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/devices/${id}`)
  },

  async discover(data: DiscoverRequest): Promise<DiscoverResult> {
    const res = await api.post<APIResponse<DiscoverResult>>('/devices/discover', data)
    return res.data.data!
  },

  async importDiscovered(ips: string[], snmpCommunity?: string, snmpVersion?: string): Promise<{ created: Device[]; skipped: string[] }> {
    const res = await api.post<APIResponse<{ created: Device[]; skipped: string[] }>>('/devices/discover/import', {
      ips,
      snmp_community: snmpCommunity ?? 'public',
      snmp_version: snmpVersion ?? 'v2c',
    })
    return res.data.data!
  },
}
