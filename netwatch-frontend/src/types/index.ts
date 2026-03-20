// ============================================================
// NetWatch - TypeScript Interfaces
// ============================================================

export type UserRole = 'admin' | 'operator' | 'viewer'
export type DeviceType = 'mikrotik' | 'cisco' | 'huawei' | 'juniper' | 'ubiquiti' | 'generic'
export type DeviceStatus = 'online' | 'offline' | 'warning' | 'unknown'
export type SNMPVersion = 'v1' | 'v2c' | 'v3'
export type AlertSeverity = 'critical' | 'high' | 'medium' | 'low' | 'info'
export type AlertStatus = 'active' | 'resolved' | 'acknowledged'
export type AlertCondition = 'gt' | 'gte' | 'lt' | 'lte' | 'eq' | 'neq' | 'down'
export type MetricType =
  | 'cpu'
  | 'memory'
  | 'temperature'
  | 'uptime'
  | 'if_in_octets'
  | 'if_out_octets'
  | 'if_in_errors'
  | 'if_out_errors'
  | 'ping_rtt'
  | 'custom'

// ---- Auth ----

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  active: boolean
  last_login_at: string | null
  created_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: number
}

export interface AuthResponse {
  tokens: TokenPair
  user: User
}

// ---- Devices ----

export interface Device {
  id: string
  name: string
  hostname: string | null
  ip_address: string
  type: DeviceType
  status: DeviceStatus
  description: string | null
  location: string | null
  snmp_version: SNMPVersion
  snmp_community: string
  snmp_port: number
  snmpv3_username?: string
  snmpv3_auth_proto?: string
  snmpv3_priv_proto?: string
  polling_interval: number
  enabled: boolean
  sys_name?: string
  sys_descr?: string
  sys_oid?: string
  sys_contact?: string
  tags: string[]
  last_seen_at: string | null
  last_polled_at: string | null
  created_at: string
  updated_at: string
}

export interface CreateDeviceRequest {
  name: string
  ip_address: string
  hostname?: string
  type?: DeviceType
  description?: string
  location?: string
  snmp_version?: SNMPVersion
  snmp_community?: string
  snmp_port?: number
  polling_interval?: number
}

export interface UpdateDeviceRequest extends Partial<CreateDeviceRequest> {
  enabled?: boolean
}

export interface DiscoverRequest {
  network: string
  snmp_version?: SNMPVersion
  snmp_community?: string
  snmp_port?: number
  timeout_seconds?: number
}

export interface DiscoveredDevice {
  ip_address: string
  sys_name: string
  sys_descr: string
  vendor: string
  type: DeviceType
  is_new: boolean
}

export interface DiscoverResult {
  network: string
  found: number
  devices: DiscoveredDevice[]
}

// ---- Metrics ----

export interface Metric {
  id: string
  device_id: string
  type: MetricType
  interface_index: number
  oid: string
  value: number
  unit: string
  collected_at: string
}

export interface MetricSummary {
  device_id: string
  type: MetricType
  min: number
  max: number
  avg: number
  last: number
  last_at: string
}

export interface MetricHistoryQuery {
  type?: MetricType
  from?: string
  to?: string
  limit?: number
}

// ---- Alerts ----

export interface NotificationConfig {
  type: 'email' | 'telegram' | 'webhook'
  config: Record<string, string>
}

export interface AlertRule {
  id: string
  name: string
  description: string | null
  device_id: string | null
  metric_type: MetricType
  condition: AlertCondition
  threshold: number
  severity: AlertSeverity
  enabled: boolean
  notifications: NotificationConfig[]
  created_at: string
  updated_at: string
}

export interface CreateAlertRuleRequest {
  name: string
  description?: string
  device_id?: string | null
  metric_type: MetricType
  condition: AlertCondition
  threshold: number
  severity: AlertSeverity
  notifications?: NotificationConfig[]
}

export interface AlertEvent {
  id: string
  rule_id: string
  device_id: string
  status: AlertStatus
  severity: AlertSeverity
  value: number | null
  message: string
  triggered_at: string
  resolved_at: string | null
  acked_at: string | null
  acked_by_id: string | null
  duration_seconds: number | null
  rule?: AlertRule
  device?: Device
}

export interface DowntimeEvent {
  id: string
  device_id: string
  started_at: string
  ended_at: string | null
  duration_seconds: number | null
  cause: string | null
  notes: string | null
}

// ---- Dashboard ----

export interface DashboardStats {
  total_devices: number
  online: number
  offline: number
  warning: number
  unknown: number
  active_alerts: number
}

// ---- API ----

export interface APIResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
  meta?: PaginationMeta
}

export interface PaginationMeta {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface PaginatedResult<T> {
  data: T[]
  meta: PaginationMeta
}

// ---- Maps ----

export interface NetworkMap {
  id: string
  name: string
  description: string | null
  layout: MapLayout
  created_by: string
  created_at: string
  updated_at: string
}

export interface MapLayout {
  nodes: MapNode[]
  edges: MapEdge[]
}

export interface MapNode {
  id: string
  type: 'device' | 'cloud' | 'switch' | 'router'
  position: { x: number; y: number }
  data: {
    device_id?: string
    label: string
    status?: DeviceStatus
    ip_address?: string
  }
}

export interface MapEdge {
  id: string
  source: string
  target: string
  label?: string
}

// ---- Electron ----

declare global {
  interface Window {
    electronAPI?: {
      getVersion: () => Promise<string>
      toggleTheme: () => Promise<string>
      store: {
        get:    (key: string) => Promise<unknown>
        set:    (key: string, value: unknown) => Promise<void>
        delete: (key: string) => Promise<void>
      }
    }
  }
}
