import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Bell, CheckCheck } from 'lucide-react'
import { Button } from '@/components/common/Button'
import { Select } from '@/components/common/Input'
import { Table, Pagination } from '@/components/common/Table'
import { SeverityBadge, Badge } from '@/components/common/Badge'
import { alertService } from '@/services/alert.service'
import { formatDistanceToNow } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import type { AlertEvent } from '@/types'

const STATUS_LABELS: Record<string, string> = {
  active: 'Ativo', resolved: 'Resolvido', acknowledged: 'Reconhecido',
}
const STATUS_VARIANT = { active: 'red', resolved: 'green', acknowledged: 'gray' } as const

export default function Alerts() {
  const queryClient = useQueryClient()
  const [page,     setPage]     = useState(1)
  const [status,   setStatus]   = useState('active')
  const [severity, setSeverity] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['alert-events', page, status, severity],
    queryFn:  () => alertService.listEvents({ page, limit: 20, status, severity }),
    refetchInterval: 15_000,
  })

  const ackMutation = useMutation({
    mutationFn: (id: string) => alertService.acknowledgeEvent(id),
    onSuccess:  () => queryClient.invalidateQueries({ queryKey: ['alert-events'] }),
  })

  const columns = [
    {
      key: 'severity', header: 'Severidade',
      render: (ev: AlertEvent) => <SeverityBadge severity={ev.severity} />,
    },
    {
      key: 'message', header: 'Mensagem',
      render: (ev: AlertEvent) => <span className="text-slate-200">{ev.message || '—'}</span>,
    },
    {
      key: 'status', header: 'Status',
      render: (ev: AlertEvent) => (
        <Badge variant={STATUS_VARIANT[ev.status] ?? 'gray'}>
          {STATUS_LABELS[ev.status] ?? ev.status}
        </Badge>
      ),
    },
    {
      key: 'when', header: 'Quando',
      render: (ev: AlertEvent) => (
        <span className="text-slate-400 text-sm">
          {formatDistanceToNow(new Date(ev.triggered_at), { locale: ptBR, addSuffix: true })}
        </span>
      ),
    },
    {
      key: 'duration', header: 'Duração',
      render: (ev: AlertEvent) => ev.duration_seconds != null
        ? <span className="text-slate-400 text-sm">{Math.round(ev.duration_seconds / 60)}min</span>
        : <span className="text-slate-600">—</span>,
    },
    {
      key: 'actions', header: '',
      className: 'w-32',
      render: (ev: AlertEvent) => ev.status === 'active' ? (
        <Button
          size="xs"
          variant="outline"
          loading={ackMutation.isPending}
          onClick={() => ackMutation.mutate(ev.id)}
          leftIcon={<CheckCheck className="w-3.5 h-3.5" />}
        >
          ACK
        </Button>
      ) : null,
    },
  ]

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-slate-100">Alertas</h1>
          <p className="text-sm text-slate-400 mt-0.5">{data?.meta?.total ?? 0} eventos</p>
        </div>
        <div className="flex items-center gap-2">
          <Bell className="w-5 h-5 text-slate-400" />
        </div>
      </div>

      {/* Filtros */}
      <div className="flex gap-3">
        <Select value={status} onChange={(e) => { setStatus(e.target.value); setPage(1) }} className="w-44">
          <option value="">Todos os status</option>
          <option value="active">Ativos</option>
          <option value="resolved">Resolvidos</option>
          <option value="acknowledged">Reconhecidos</option>
        </Select>
        <Select value={severity} onChange={(e) => { setSeverity(e.target.value); setPage(1) }} className="w-44">
          <option value="">Todas as severidades</option>
          <option value="critical">Crítico</option>
          <option value="high">Alto</option>
          <option value="medium">Médio</option>
          <option value="low">Baixo</option>
          <option value="info">Info</option>
        </Select>
      </div>

      <Table columns={columns} data={data?.data ?? []} keyFn={(ev) => ev.id} loading={isLoading} />
      {data?.meta && data.meta.total_pages > 1 && (
        <Pagination meta={data.meta} onChange={setPage} />
      )}
    </div>
  )
}
