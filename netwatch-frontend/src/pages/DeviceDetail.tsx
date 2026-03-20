import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft, Cpu, MemoryStick, Thermometer, Activity } from 'lucide-react'
import { StatusBadge } from '@/components/common/Badge'
import { Card } from '@/components/common/Card'
import { deviceService } from '@/services/device.service'
import { metricService } from '@/services/metric.service'
import { formatDistanceToNow, subHours } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import {
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from 'recharts'
import { format } from 'date-fns'

export default function DeviceDetail() {
  const { id } = useParams<{ id: string }>()

  const { data: device } = useQuery({
    queryKey: ['device', id],
    queryFn:  () => deviceService.getById(id!),
    enabled:  !!id,
  })

  const from = subHours(new Date(), 24).toISOString()
  const to   = new Date().toISOString()

  const { data: cpuMetrics } = useQuery({
    queryKey: ['metrics', id, 'cpu'],
    queryFn:  () => metricService.getHistory(id!, { type: 'cpu', from, to }),
    enabled:  !!id,
    refetchInterval: 60_000,
  })

  const { data: memMetrics } = useQuery({
    queryKey: ['metrics', id, 'memory'],
    queryFn:  () => metricService.getHistory(id!, { type: 'memory', from, to }),
    enabled:  !!id,
    refetchInterval: 60_000,
  })

  if (!device) return (
    <div className="flex items-center justify-center h-64 text-slate-500">
      <div className="w-6 h-6 border-2 border-brand-500 border-t-transparent rounded-full animate-spin mr-2" />
      Carregando...
    </div>
  )

  const cpuData = (cpuMetrics ?? []).map((m) => ({
    time:  format(new Date(m.collected_at), 'HH:mm'),
    value: Math.round(m.value),
  })).reverse()

  const memData = (memMetrics ?? []).map((m) => ({
    time:  format(new Date(m.collected_at), 'HH:mm'),
    value: Math.round(m.value / 1024 / 1024),
  })).reverse()

  const chartProps = {
    style: { fontSize: 11, fontFamily: 'inherit' },
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-3">
        <Link to="/devices" className="flex items-center gap-1.5 text-slate-400 hover:text-slate-200 text-sm transition-colors">
          <ArrowLeft className="w-4 h-4" /> Dispositivos
        </Link>
        <span className="text-slate-600">/</span>
        <span className="text-slate-300 text-sm">{device.name}</span>
      </div>

      {/* Header do dispositivo */}
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center gap-3 mb-1">
            <h1 className="text-xl font-bold text-slate-100">{device.name}</h1>
            <StatusBadge status={device.status} />
          </div>
          <p className="text-sm text-slate-400 font-mono">{device.ip_address}</p>
          {device.last_seen_at && (
            <p className="text-xs text-slate-500 mt-1">
              Visto {formatDistanceToNow(new Date(device.last_seen_at), { locale: ptBR, addSuffix: true })}
            </p>
          )}
        </div>
      </div>

      {/* Info cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <InfoItem label="Tipo" value={device.type} />
        <InfoItem label="SNMP" value={`${device.snmp_version} / ${device.snmp_community}`} />
        <InfoItem label="Polling" value={`${device.polling_interval}s`} />
        <InfoItem label="Localização" value={device.location ?? '—'} />
      </div>

      {/* sysInfo */}
      {(device.sys_name || device.sys_descr) && (
        <Card>
          <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">System Info (SNMP)</p>
          <div className="grid grid-cols-2 gap-3 text-sm">
            {device.sys_name    && <InfoItem label="sysName"    value={device.sys_name} />}
            {device.sys_contact && <InfoItem label="sysContact" value={device.sys_contact} />}
            {device.sys_descr   && <InfoItem label="sysDescr"   value={device.sys_descr} />}
            {device.sys_oid     && <InfoItem label="sysOID"     value={device.sys_oid} mono />}
          </div>
        </Card>
      )}

      {/* Gráficos */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <MetricChart
          title="CPU (%)"
          icon={<Cpu className="w-4 h-4 text-blue-400" />}
          data={cpuData}
          unit="%"
          color="#3b82f6"
          domain={[0, 100]}
        />
        <MetricChart
          title="Memória (MB)"
          icon={<MemoryStick className="w-4 h-4 text-purple-400" />}
          data={memData}
          unit=" MB"
          color="#a855f7"
        />
      </div>
    </div>
  )
}

function InfoItem({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div>
      <p className="text-xs text-slate-500">{label}</p>
      <p className={`text-sm text-slate-200 mt-0.5 ${mono ? 'font-mono' : ''}`}>{value}</p>
    </div>
  )
}

interface MetricChartProps {
  title:  string
  icon:   React.ReactNode
  data:   { time: string; value: number }[]
  unit:   string
  color:  string
  domain?: [number, number]
}

function MetricChart({ title, icon, data, unit, color, domain }: MetricChartProps) {
  return (
    <Card>
      <div className="flex items-center gap-2 mb-4">
        {icon}
        <span className="text-sm font-semibold text-slate-300">{title}</span>
        {data.length > 0 && (
          <span className="ml-auto text-lg font-bold text-slate-100">
            {data[data.length - 1]?.value}{unit}
          </span>
        )}
      </div>
      {data.length === 0 ? (
        <div className="h-40 flex items-center justify-center text-slate-500 text-sm">
          Sem dados disponíveis
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={160}>
          <LineChart data={data} margin={{ top: 2, right: 2, left: -20, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="time" tick={{ fill: '#64748b', fontSize: 10 }} interval="preserveStartEnd" />
            <YAxis tick={{ fill: '#64748b', fontSize: 10 }} domain={domain} />
            <Tooltip
              contentStyle={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 6 }}
              itemStyle={{ color: '#cbd5e1' }}
              formatter={(v: number) => [`${v}${unit}`, title]}
            />
            <Line type="monotone" dataKey="value" stroke={color} strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      )}
    </Card>
  )
}
