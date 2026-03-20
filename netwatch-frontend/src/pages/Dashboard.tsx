import { useQuery } from '@tanstack/react-query'
import { Monitor, Wifi, WifiOff, AlertTriangle, Bell, RefreshCw } from 'lucide-react'
import { StatCard } from '@/components/common/Card'
import { StatusBadge, SeverityBadge } from '@/components/common/Badge'
import { dashboardService } from '@/services/dashboard.service'
import { deviceService } from '@/services/device.service'
import { alertService } from '@/services/alert.service'
import { formatDistanceToNow } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import {
  PieChart, Pie, Cell, Tooltip, ResponsiveContainer,
} from 'recharts'
import { Link } from 'react-router-dom'

const STATUS_COLORS = {
  online:  '#22c55e',
  offline: '#ef4444',
  warning: '#f59e0b',
  unknown: '#6b7280',
}

export default function Dashboard() {
  const { data: stats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn:  dashboardService.getStats,
    refetchInterval: 30_000,
  })

  const { data: recentDevices } = useQuery({
    queryKey: ['devices-recent'],
    queryFn:  () => deviceService.list({ limit: 5, status: 'offline' }),
    refetchInterval: 30_000,
  })

  const { data: activeAlerts } = useQuery({
    queryKey: ['alerts-active'],
    queryFn:  () => alertService.listEvents({ limit: 5, status: 'active' }),
    refetchInterval: 30_000,
  })

  const pieData = stats ? [
    { name: 'Online',  value: stats.online,  color: STATUS_COLORS.online },
    { name: 'Offline', value: stats.offline, color: STATUS_COLORS.offline },
    { name: 'Warning', value: stats.warning, color: STATUS_COLORS.warning },
    { name: 'Unknown', value: stats.unknown, color: STATUS_COLORS.unknown },
  ].filter((d) => d.value > 0) : []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-slate-100">Dashboard</h1>
          <p className="text-sm text-slate-400 mt-0.5">Visão geral do ambiente monitorado</p>
        </div>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total de Dispositivos"
          value={stats?.total_devices ?? '—'}
          icon={<Monitor className="w-5 h-5" />}
          color="blue"
        />
        <StatCard
          title="Online"
          value={stats?.online ?? '—'}
          icon={<Wifi className="w-5 h-5" />}
          color="green"
        />
        <StatCard
          title="Offline"
          value={stats?.offline ?? '—'}
          icon={<WifiOff className="w-5 h-5" />}
          color="red"
        />
        <StatCard
          title="Alertas Ativos"
          value={stats?.active_alerts ?? '—'}
          icon={<Bell className="w-5 h-5" />}
          color="yellow"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Gráfico de status */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-5">
          <h2 className="text-sm font-semibold text-slate-300 mb-4">Status dos Dispositivos</h2>
          {pieData.length > 0 ? (
            <div className="flex items-center gap-6">
              <ResponsiveContainer width={140} height={140}>
                <PieChart>
                  <Pie data={pieData} cx="50%" cy="50%" innerRadius={40} outerRadius={65} dataKey="value" strokeWidth={0}>
                    {pieData.map((entry, i) => (
                      <Cell key={i} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 8 }}
                    itemStyle={{ color: '#cbd5e1' }}
                  />
                </PieChart>
              </ResponsiveContainer>
              <div className="space-y-2">
                {pieData.map((d) => (
                  <div key={d.name} className="flex items-center gap-2 text-sm">
                    <span className="w-2.5 h-2.5 rounded-full" style={{ background: d.color }} />
                    <span className="text-slate-400">{d.name}</span>
                    <span className="text-slate-200 font-medium ml-auto pl-4">{d.value}</span>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="h-40 flex items-center justify-center text-slate-500">
              Nenhum dispositivo cadastrado
            </div>
          )}
        </div>

        {/* Dispositivos offline recentes */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-5">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-semibold text-slate-300">Dispositivos Offline</h2>
            <Link to="/devices?status=offline" className="text-xs text-brand-400 hover:text-brand-300">
              Ver todos
            </Link>
          </div>
          <div className="space-y-2">
            {recentDevices?.data.length === 0 && (
              <p className="text-sm text-slate-500 py-4 text-center">Nenhum dispositivo offline</p>
            )}
            {recentDevices?.data.map((d) => (
              <Link
                key={d.id}
                to={`/devices/${d.id}`}
                className="flex items-center justify-between p-3 rounded-lg hover:bg-slate-700/50 transition-colors"
              >
                <div>
                  <p className="text-sm font-medium text-slate-200">{d.name}</p>
                  <p className="text-xs text-slate-500 font-mono">{d.ip_address}</p>
                </div>
                <div className="flex items-center gap-2">
                  {d.last_seen_at && (
                    <span className="text-xs text-slate-500">
                      {formatDistanceToNow(new Date(d.last_seen_at), { locale: ptBR, addSuffix: true })}
                    </span>
                  )}
                  <StatusBadge status={d.status} />
                </div>
              </Link>
            ))}
          </div>
        </div>
      </div>

      {/* Alertas ativos */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-5">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-slate-300">Alertas Ativos</h2>
          <Link to="/alerts" className="text-xs text-brand-400 hover:text-brand-300">Ver todos</Link>
        </div>
        {activeAlerts?.data.length === 0 ? (
          <div className="py-8 text-center">
            <div className="w-10 h-10 bg-green-400/10 rounded-full flex items-center justify-center mx-auto mb-2">
              <RefreshCw className="w-5 h-5 text-green-400" />
            </div>
            <p className="text-sm text-slate-400">Sem alertas ativos</p>
          </div>
        ) : (
          <div className="space-y-2">
            {activeAlerts?.data.map((ev) => (
              <div key={ev.id} className="flex items-center justify-between p-3 rounded-lg bg-slate-700/30">
                <div className="flex items-center gap-3">
                  <AlertTriangle className="w-4 h-4 text-yellow-400 shrink-0" />
                  <div>
                    <p className="text-sm text-slate-200">{ev.message}</p>
                    <p className="text-xs text-slate-500">
                      {formatDistanceToNow(new Date(ev.triggered_at), { locale: ptBR, addSuffix: true })}
                    </p>
                  </div>
                </div>
                <SeverityBadge severity={ev.severity} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
