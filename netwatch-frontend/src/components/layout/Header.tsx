import { useAuthStore } from '@/stores/authStore'
import { useUIStore } from '@/stores/uiStore'
import { LogOut, Sun, Moon, User } from 'lucide-react'
import { Badge } from '@/components/common/Badge'

export function Header() {
  const { user, logout } = useAuthStore()
  const { theme, toggleTheme } = useUIStore()

  const roleVariant = { admin: 'red', operator: 'yellow', viewer: 'gray' } as const

  return (
    <header className="h-14 bg-slate-900 border-b border-slate-700/50 flex items-center justify-between px-5 shrink-0">
      {/* Breadcrumb / título pode ser injetado pelo contexto */}
      <div />

      {/* Ações direita */}
      <div className="flex items-center gap-3">
        {/* Toggle tema */}
        <button
          onClick={toggleTheme}
          className="p-2 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-700/50 transition-colors"
          title={theme === 'dark' ? 'Mudar para Light' : 'Mudar para Dark'}
        >
          {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
        </button>

        {/* User info */}
        {user && (
          <div className="flex items-center gap-2.5 pl-3 border-l border-slate-700">
            <div className="w-7 h-7 bg-brand-700 rounded-full flex items-center justify-center">
              <User className="w-4 h-4 text-brand-200" />
            </div>
            <div className="hidden sm:block">
              <p className="text-sm font-medium text-slate-200 leading-none">{user.name}</p>
              <div className="mt-0.5">
                <Badge variant={roleVariant[user.role]}>{user.role}</Badge>
              </div>
            </div>
          </div>
        )}

        {/* Logout */}
        <button
          onClick={() => logout()}
          className="p-2 rounded-lg text-slate-400 hover:text-red-400 hover:bg-red-400/10 transition-colors"
          title="Sair"
        >
          <LogOut className="w-4 h-4" />
        </button>
      </div>
    </header>
  )
}
