import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { Plus, Search, RefreshCw, Trash2, Edit, Eye } from 'lucide-react'
import { Button } from '@/components/common/Button'
import { Input, Select } from '@/components/common/Input'
import { Table, Pagination } from '@/components/common/Table'
import { StatusBadge } from '@/components/common/Badge'
import { ConfirmModal } from '@/components/common/Modal'
import { deviceService } from '@/services/device.service'
import { useIsOperator } from '@/stores/authStore'
import { formatDistanceToNow } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import type { Device } from '@/types'
import { DeviceModal } from '@/components/devices/DeviceModal'

const DEVICE_TYPE_LABELS: Record<string, string> = {
  mikrotik: 'MikroTik',
  cisco:    'Cisco',
  huawei:   'Huawei',
  juniper:  'Juniper',
  ubiquiti: 'Ubiquiti',
  generic:  'Genérico',
}

export default function Devices() {
  const queryClient = useQueryClient()
  const isOperator  = useIsOperator()

  const [page,   setPage]   = useState(1)
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState('')
  const [type,   setType]   = useState('')

  const [modalOpen,   setModalOpen]   = useState(false)
  const [editDevice,  setEditDevice]  = useState<Device | null>(null)
  const [deleteId,    setDeleteId]    = useState<string | null>(null)

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['devices', page, search, status, type],
    queryFn:  () => deviceService.list({ page, limit: 20, search, status, type }),
    refetchInterval: 60_000,
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deviceService.delete(id),
    onSuccess:  () => {
      queryClient.invalidateQueries({ queryKey: ['devices'] })
      setDeleteId(null)
    },
  })

  const columns = [
    {
      key: 'name', header: 'Nome',
      render: (d: Device) => (
        <div>
          <Link to={`/devices/${d.id}`} className="font-medium text-slate-200 hover:text-brand-400 transition-colors">
            {d.name}
          </Link>
          {d.location && <p className="text-xs text-slate-500">{d.location}</p>}
        </div>
      ),
    },
    {
      key: 'ip', header: 'IP',
      render: (d: Device) => <span className="font-mono text-sm">{d.ip_address}</span>,
    },
    {
      key: 'type', header: 'Tipo',
      render: (d: Device) => <span className="text-slate-400">{DEVICE_TYPE_LABELS[d.type] ?? d.type}</span>,
    },
    {
      key: 'status', header: 'Status',
      render: (d: Device) => <StatusBadge status={d.status} />,
    },
    {
      key: 'seen', header: 'Último Contato',
      render: (d: Device) => d.last_seen_at
        ? <span className="text-slate-400 text-sm">{formatDistanceToNow(new Date(d.last_seen_at), { locale: ptBR, addSuffix: true })}</span>
        : <span className="text-slate-600">—</span>,
    },
    {
      key: 'actions', header: '',
      className: 'w-28',
      render: (d: Device) => (
        <div className="flex items-center gap-1">
          <Link to={`/devices/${d.id}`} className="p-1.5 text-slate-400 hover:text-slate-200 hover:bg-slate-700 rounded transition-colors" title="Detalhes">
            <Eye className="w-4 h-4" />
          </Link>
          {isOperator && (
            <>
              <button
                onClick={() => { setEditDevice(d); setModalOpen(true) }}
                className="p-1.5 text-slate-400 hover:text-slate-200 hover:bg-slate-700 rounded transition-colors"
                title="Editar"
              >
                <Edit className="w-4 h-4" />
              </button>
              <button
                onClick={() => setDeleteId(d.id)}
                className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-red-400/10 rounded transition-colors"
                title="Excluir"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </>
          )}
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-5">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-slate-100">Dispositivos</h1>
          <p className="text-sm text-slate-400 mt-0.5">{data?.meta?.total ?? 0} dispositivos cadastrados</p>
        </div>
        <div className="flex gap-2">
          <Button variant="ghost" size="sm" onClick={() => refetch()} leftIcon={<RefreshCw className="w-4 h-4" />}>
            Atualizar
          </Button>
          {isOperator && (
            <Button size="sm" onClick={() => { setEditDevice(null); setModalOpen(true) }} leftIcon={<Plus className="w-4 h-4" />}>
              Adicionar
            </Button>
          )}
        </div>
      </div>

      {/* Filtros */}
      <div className="flex flex-wrap gap-3">
        <div className="flex-1 min-w-48">
          <Input
            placeholder="Buscar por nome ou IP..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1) }}
            leftIcon={<Search className="w-4 h-4" />}
          />
        </div>
        <div className="w-40">
          <Select value={status} onChange={(e) => { setStatus(e.target.value); setPage(1) }}>
            <option value="">Todos os status</option>
            <option value="online">Online</option>
            <option value="offline">Offline</option>
            <option value="warning">Warning</option>
            <option value="unknown">Unknown</option>
          </Select>
        </div>
        <div className="w-40">
          <Select value={type} onChange={(e) => { setType(e.target.value); setPage(1) }}>
            <option value="">Todos os tipos</option>
            <option value="mikrotik">MikroTik</option>
            <option value="cisco">Cisco</option>
            <option value="huawei">Huawei</option>
            <option value="ubiquiti">Ubiquiti</option>
            <option value="generic">Genérico</option>
          </Select>
        </div>
      </div>

      {/* Tabela */}
      <Table columns={columns} data={data?.data ?? []} keyFn={(d) => d.id} loading={isLoading} />
      {data?.meta && data.meta.total_pages > 1 && (
        <Pagination meta={data.meta} onChange={setPage} />
      )}

      {/* Modal criar/editar */}
      <DeviceModal
        open={modalOpen}
        device={editDevice}
        onClose={() => { setModalOpen(false); setEditDevice(null) }}
        onSaved={() => queryClient.invalidateQueries({ queryKey: ['devices'] })}
      />

      {/* Confirm delete */}
      <ConfirmModal
        open={!!deleteId}
        onClose={() => setDeleteId(null)}
        onConfirm={() => deleteId && deleteMutation.mutate(deleteId)}
        loading={deleteMutation.isPending}
        title="Remover Dispositivo"
        message="Tem certeza que deseja remover este dispositivo? Esta ação não pode ser desfeita."
        confirmLabel="Remover"
      />
    </div>
  )
}
