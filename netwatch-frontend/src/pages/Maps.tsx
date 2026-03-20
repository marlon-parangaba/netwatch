import { useCallback, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  ReactFlow,
  Node, Edge, NodeChange, EdgeChange, Connection,
  addEdge, applyNodeChanges, applyEdgeChanges,
  MiniMap, Controls, Background, BackgroundVariant,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { deviceService } from '@/services/device.service'
import type { Device } from '@/types'

const STATUS_COLOR: Record<string, string> = {
  online:  '#22c55e',
  offline: '#ef4444',
  warning: '#f59e0b',
  unknown: '#6b7280',
}

function buildNodesFromDevices(devices: Device[]): Node[] {
  return devices.map((d, i) => ({
    id:       d.id,
    type:     'default',
    position: { x: (i % 5) * 200, y: Math.floor(i / 5) * 150 },
    data: {
      label: (
        <div className="text-center">
          <div
            className="w-3 h-3 rounded-full mx-auto mb-1"
            style={{ background: STATUS_COLOR[d.status] ?? '#6b7280' }}
          />
          <p className="text-xs font-medium text-slate-200 truncate max-w-24">{d.name}</p>
          <p className="text-xs text-slate-500 font-mono">{d.ip_address}</p>
        </div>
      ),
    },
    style: {
      background: '#1e293b',
      border:     `1px solid ${STATUS_COLOR[d.status] ?? '#334155'}`,
      borderRadius: 8,
      width:      120,
      padding:    8,
    },
  }))
}

export default function Maps() {
  const { data: devicesResult } = useQuery({
    queryKey: ['devices-map'],
    queryFn:  () => deviceService.list({ limit: 200 }),
  })

  const [nodes, setNodes] = useState<Node[]>([])
  const [edges, setEdges] = useState<Edge[]>([])
  const [loaded, setLoaded] = useState(false)

  // Carregar dispositivos no mapa na primeira vez
  if (devicesResult && !loaded) {
    setNodes(buildNodesFromDevices(devicesResult.data))
    setLoaded(true)
  }

  const onNodesChange = useCallback(
    (changes: NodeChange[]) => setNodes((nds) => applyNodeChanges(changes, nds)),
    []
  )
  const onEdgesChange = useCallback(
    (changes: EdgeChange[]) => setEdges((eds) => applyEdgeChanges(changes, eds)),
    []
  )
  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    []
  )

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-xl font-bold text-slate-100">Mapa de Rede</h1>
        <p className="text-sm text-slate-400 mt-0.5">
          Visualização topológica — arraste dispositivos, conecte com linhas
        </p>
      </div>

      <div className="h-[calc(100vh-220px)] rounded-xl overflow-hidden border border-slate-700">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          fitView
          colorMode="dark"
        >
          <Background variant={BackgroundVariant.Dots} gap={24} color="#334155" />
          <MiniMap
            style={{ background: '#0f172a', border: '1px solid #334155' }}
            nodeColor={(n) => {
              const status = (n.data as { status?: string })?.status ?? 'unknown'
              return STATUS_COLOR[status] ?? '#6b7280'
            }}
          />
          <Controls style={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 8 }} />
        </ReactFlow>
      </div>
    </div>
  )
}
