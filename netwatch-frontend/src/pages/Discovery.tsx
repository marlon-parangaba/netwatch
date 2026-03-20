import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { Radar, Download, CheckSquare, Square } from 'lucide-react'
import { Button } from '@/components/common/Button'
import { Input, Select } from '@/components/common/Input'
import { Card } from '@/components/common/Card'
import { StatusBadge } from '@/components/common/Badge'
import { deviceService } from '@/services/device.service'
import type { DiscoveredDevice } from '@/types'

export default function Discovery() {
  const [network,    setNetwork]    = useState('192.168.1.0/24')
  const [community,  setCommunity]  = useState('public')
  const [version,    setVersion]    = useState('v2c')
  const [results,    setResults]    = useState<DiscoveredDevice[] | null>(null)
  const [selected,   setSelected]   = useState<Set<string>>(new Set())

  const discoverMutation = useMutation({
    mutationFn: () => deviceService.discover({ network, snmp_community: community, snmp_version: version as 'v2c', timeout_seconds: 5 }),
    onSuccess:  (data) => {
      setResults(data.devices)
      setSelected(new Set(data.devices.filter((d) => d.is_new).map((d) => d.ip_address)))
    },
  })

  const importMutation = useMutation({
    mutationFn: () => deviceService.importDiscovered([...selected], community, version),
    onSuccess:  () => {
      alert(`${selected.size} dispositivos importados com sucesso!`)
      setSelected(new Set())
    },
  })

  const toggleSelect = (ip: string) => {
    setSelected((prev) => {
      const next = new Set(prev)
      if (next.has(ip)) next.delete(ip)
      else next.add(ip)
      return next
    })
  }

  const toggleAll = () => {
    if (!results) return
    const newOnly = results.filter((d) => d.is_new).map((d) => d.ip_address)
    if (selected.size === newOnly.length) setSelected(new Set())
    else setSelected(new Set(newOnly))
  }

  return (
    <div className="space-y-6 max-w-4xl">
      <div>
        <h1 className="text-xl font-bold text-slate-100">Discovery de Rede</h1>
        <p className="text-sm text-slate-400 mt-0.5">Varre uma rede e detecta dispositivos via SNMP</p>
      </div>

      {/* Configuração */}
      <Card>
        <div className="space-y-4">
          <p className="text-sm font-semibold text-slate-300">Parâmetros de Varredura</p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="sm:col-span-1">
              <Input
                label="Rede CIDR"
                value={network}
                onChange={(e) => setNetwork(e.target.value)}
                placeholder="192.168.1.0/24"
              />
            </div>
            <Input
              label="SNMP Community"
              value={community}
              onChange={(e) => setCommunity(e.target.value)}
              placeholder="public"
            />
            <Select label="Versão SNMP" value={version} onChange={(e) => setVersion(e.target.value)}>
              <option value="v1">v1</option>
              <option value="v2c">v2c</option>
              <option value="v3">v3</option>
            </Select>
          </div>
          <Button
            onClick={() => discoverMutation.mutate()}
            loading={discoverMutation.isPending}
            leftIcon={<Radar className="w-4 h-4" />}
          >
            {discoverMutation.isPending ? 'Varrendo...' : 'Iniciar Varredura'}
          </Button>
          {discoverMutation.isPending && (
            <div className="flex items-center gap-3 mt-2">
              <div className="flex-1 h-2 bg-slate-700 rounded-full overflow-hidden">
                <div className="h-full bg-brand-500 rounded-full animate-pulse w-1/2" />
              </div>
              <span className="text-xs text-slate-400">Aguarde...</span>
            </div>
          )}
        </div>
      </Card>

      {/* Resultados */}
      {results !== null && (
        <Card>
          <div className="flex items-center justify-between mb-4">
            <p className="text-sm font-semibold text-slate-300">
              {results.length} dispositivo{results.length !== 1 ? 's' : ''} encontrado{results.length !== 1 ? 's' : ''}
            </p>
            {results.some((d) => d.is_new) && (
              <div className="flex items-center gap-3">
                <button onClick={toggleAll} className="text-xs text-slate-400 hover:text-slate-200 flex items-center gap-1">
                  {selected.size > 0 ? <CheckSquare className="w-4 h-4" /> : <Square className="w-4 h-4" />}
                  Selecionar novos
                </button>
                <Button
                  size="sm"
                  onClick={() => importMutation.mutate()}
                  loading={importMutation.isPending}
                  disabled={selected.size === 0}
                  leftIcon={<Download className="w-4 h-4" />}
                >
                  Importar ({selected.size})
                </Button>
              </div>
            )}
          </div>

          {results.length === 0 ? (
            <div className="py-12 text-center text-slate-500">
              Nenhum dispositivo respondeu ao SNMP na rede {network}
            </div>
          ) : (
            <div className="space-y-2">
              {results.map((d) => (
                <div
                  key={d.ip_address}
                  className={`flex items-center gap-3 p-3 rounded-lg border transition-colors cursor-pointer ${
                    selected.has(d.ip_address)
                      ? 'border-brand-500/50 bg-brand-500/10'
                      : 'border-slate-700 hover:border-slate-600'
                  }`}
                  onClick={() => d.is_new && toggleSelect(d.ip_address)}
                >
                  {d.is_new ? (
                    selected.has(d.ip_address)
                      ? <CheckSquare className="w-4 h-4 text-brand-400 shrink-0" />
                      : <Square className="w-4 h-4 text-slate-500 shrink-0" />
                  ) : (
                    <div className="w-4 h-4 shrink-0" />
                  )}

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-sm text-slate-200">{d.ip_address}</span>
                      {d.sys_name && <span className="text-slate-400 text-sm truncate">({d.sys_name})</span>}
                    </div>
                    {d.sys_descr && (
                      <p className="text-xs text-slate-500 truncate mt-0.5">{d.sys_descr}</p>
                    )}
                  </div>

                  <div className="flex items-center gap-2 shrink-0">
                    <span className="text-xs text-slate-400 capitalize">{d.vendor}</span>
                    {!d.is_new && (
                      <StatusBadge status="online" />
                    )}
                    {!d.is_new && (
                      <span className="text-xs text-slate-500 bg-slate-700 px-2 py-0.5 rounded">Já cadastrado</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}
    </div>
  )
}
