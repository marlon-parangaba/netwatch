import { useState, useEffect } from 'react'
import { useMutation } from '@tanstack/react-query'
import { Modal } from '@/components/common/Modal'
import { Input, Select } from '@/components/common/Input'
import { Button } from '@/components/common/Button'
import { deviceService } from '@/services/device.service'
import type { Device, CreateDeviceRequest } from '@/types'

interface DeviceModalProps {
  open:     boolean
  device:   Device | null
  onClose:  () => void
  onSaved:  () => void
}

const defaultForm: CreateDeviceRequest = {
  name:             '',
  ip_address:       '',
  hostname:         '',
  type:             'generic',
  description:      '',
  location:         '',
  snmp_version:     'v2c',
  snmp_community:   'public',
  snmp_port:        161,
  polling_interval: 60,
}

export function DeviceModal({ open, device, onClose, onSaved }: DeviceModalProps) {
  const [form, setForm] = useState(defaultForm)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (device) {
      setForm({
        name:             device.name,
        ip_address:       device.ip_address,
        hostname:         device.hostname ?? '',
        type:             device.type,
        description:      device.description ?? '',
        location:         device.location ?? '',
        snmp_version:     device.snmp_version,
        snmp_community:   device.snmp_community,
        snmp_port:        device.snmp_port,
        polling_interval: device.polling_interval,
      })
    } else {
      setForm(defaultForm)
    }
    setErrors({})
  }, [device, open])

  const set = (key: keyof CreateDeviceRequest, value: unknown) =>
    setForm((f) => ({ ...f, [key]: value }))

  const validate = (): boolean => {
    const e: Record<string, string> = {}
    if (!form.name)       e.name       = 'Nome é obrigatório'
    if (!form.ip_address) e.ip_address = 'IP é obrigatório'
    setErrors(e)
    return Object.keys(e).length === 0
  }

  const createMutation = useMutation({
    mutationFn: (d: CreateDeviceRequest) => deviceService.create(d),
    onSuccess:  () => { onSaved(); onClose() },
  })

  const updateMutation = useMutation({
    mutationFn: (d: CreateDeviceRequest) => deviceService.update(device!.id, d),
    onSuccess:  () => { onSaved(); onClose() },
  })

  const handleSubmit = () => {
    if (!validate()) return
    if (device) updateMutation.mutate(form)
    else         createMutation.mutate(form)
  }

  const isPending = createMutation.isPending || updateMutation.isPending
  const apiError  = (createMutation.error || updateMutation.error) as { response?: { data?: { error?: string } } } | null

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={device ? 'Editar Dispositivo' : 'Novo Dispositivo'}
      size="lg"
      footer={
        <>
          <Button variant="ghost" onClick={onClose}>Cancelar</Button>
          <Button loading={isPending} onClick={handleSubmit}>
            {device ? 'Salvar' : 'Criar'}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        {apiError?.response?.data?.error && (
          <div className="px-4 py-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
            {apiError.response.data.error}
          </div>
        )}

        <div className="grid grid-cols-2 gap-4">
          <Input label="Nome" value={form.name} onChange={(e) => set('name', e.target.value)} error={errors.name} />
          <Input label="Endereço IP" value={form.ip_address} onChange={(e) => set('ip_address', e.target.value)} error={errors.ip_address} placeholder="192.168.1.1" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input label="Hostname" value={form.hostname} onChange={(e) => set('hostname', e.target.value)} />
          <Select label="Tipo" value={form.type} onChange={(e) => set('type', e.target.value)}>
            <option value="generic">Genérico</option>
            <option value="mikrotik">MikroTik</option>
            <option value="cisco">Cisco</option>
            <option value="huawei">Huawei</option>
            <option value="juniper">Juniper</option>
            <option value="ubiquiti">Ubiquiti</option>
          </Select>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Input label="Localização" value={form.location} onChange={(e) => set('location', e.target.value)} />
          <Input label="Descrição" value={form.description} onChange={(e) => set('description', e.target.value)} />
        </div>

        <div className="border-t border-slate-700 pt-4">
          <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Configuração SNMP</p>
          <div className="grid grid-cols-3 gap-4">
            <Select label="Versão" value={form.snmp_version} onChange={(e) => set('snmp_version', e.target.value)}>
              <option value="v1">v1</option>
              <option value="v2c">v2c</option>
              <option value="v3">v3</option>
            </Select>
            <Input label="Community" value={form.snmp_community} onChange={(e) => set('snmp_community', e.target.value)} />
            <Input label="Porta" type="number" value={form.snmp_port} onChange={(e) => set('snmp_port', Number(e.target.value))} />
          </div>
        </div>

        <div>
          <Input
            label="Intervalo de Polling (segundos)"
            type="number"
            value={form.polling_interval}
            onChange={(e) => set('polling_interval', Number(e.target.value))}
            className="w-48"
          />
        </div>
      </div>
    </Modal>
  )
}
