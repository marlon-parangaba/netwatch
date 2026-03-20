import { Card } from '@/components/common/Card'
import { useAuthStore } from '@/stores/authStore'
import { useIsAdmin } from '@/stores/authStore'
import { Shield, Info } from 'lucide-react'

export default function Settings() {
  const { user } = useAuthStore()
  const isAdmin  = useIsAdmin()

  return (
    <div className="space-y-6 max-w-2xl">
      <div>
        <h1 className="text-xl font-bold text-slate-100">Configurações</h1>
        <p className="text-sm text-slate-400 mt-0.5">Preferências e administração do sistema</p>
      </div>

      {/* Perfil */}
      <Card>
        <p className="text-sm font-semibold text-slate-300 mb-4">Meu Perfil</p>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-slate-500">Nome</p>
            <p className="text-slate-200 mt-0.5">{user?.name}</p>
          </div>
          <div>
            <p className="text-slate-500">E-mail</p>
            <p className="text-slate-200 mt-0.5">{user?.email}</p>
          </div>
          <div>
            <p className="text-slate-500">Função</p>
            <p className="text-slate-200 mt-0.5 capitalize">{user?.role}</p>
          </div>
        </div>
      </Card>

      {/* Admin only */}
      {isAdmin && (
        <Card>
          <div className="flex items-center gap-2 mb-4">
            <Shield className="w-4 h-4 text-brand-400" />
            <p className="text-sm font-semibold text-slate-300">Administração</p>
          </div>
          <p className="text-sm text-slate-400">
            Gestão completa de usuários e configurações avançadas será implementada na Fase 6 (Polimento UX).
          </p>
        </Card>
      )}

      {/* Sobre */}
      <Card>
        <div className="flex items-center gap-2 mb-3">
          <Info className="w-4 h-4 text-slate-400" />
          <p className="text-sm font-semibold text-slate-300">Sobre o NetWatch</p>
        </div>
        <div className="text-sm text-slate-400 space-y-1">
          <p>Versão: <span className="text-slate-200">0.1.0</span></p>
          <p>Backend: <span className="text-slate-200">Go + Fiber + PostgreSQL</span></p>
          <p>Frontend: <span className="text-slate-200">Electron + React + TypeScript</span></p>
        </div>
      </Card>
    </div>
  )
}
