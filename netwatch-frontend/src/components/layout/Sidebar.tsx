import { NavLink } from 'react-router-dom'
import { clsx } from 'clsx'
import {
  LayoutDashboard, Monitor, Radar, Bell, Map, Settings,
  ChevronLeft, ChevronRight, Activity,
} from 'lucide-react'
import { useUIStore } from '@/stores/uiStore'

const navItems = [
  { to: '/dashboard',  icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/devices',    icon: Monitor,          label: 'Dispositivos' },
  { to: '/discovery',  icon: Radar,            label: 'Discovery' },
  { to: '/alerts',     icon: Bell,             label: 'Alertas' },
  { to: '/maps',       icon: Map,              label: 'Mapas' },
  { to: '/settings',   icon: Settings,         label: 'Configurações' },
]

export function Sidebar() {
  const { sidebarCollapsed, toggleSidebar } = useUIStore()

  return (
    <aside className={clsx(
      'flex flex-col bg-slate-900 border-r border-slate-700/50',
      'transition-all duration-200 ease-in-out shrink-0',
      sidebarCollapsed ? 'w-16' : 'w-56'
    )}>
      {/* Logo */}
      <div className={clsx(
        'flex items-center gap-3 px-4 py-4 border-b border-slate-700/50',
        sidebarCollapsed && 'justify-center px-0'
      )}>
        <div className="w-8 h-8 bg-brand-600 rounded-lg flex items-center justify-center shrink-0">
          <Activity className="w-5 h-5 text-white" />
        </div>
        {!sidebarCollapsed && (
          <span className="font-bold text-slate-100 text-lg tracking-tight">NetWatch</span>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 py-4 space-y-1 px-2">
        {navItems.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            title={sidebarCollapsed ? label : undefined}
            className={({ isActive }) => clsx(
              'flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors text-sm font-medium',
              isActive
                ? 'bg-brand-600/20 text-brand-400'
                : 'text-slate-400 hover:text-slate-100 hover:bg-slate-700/50',
              sidebarCollapsed && 'justify-center px-0'
            )}
          >
            <Icon className="w-5 h-5 shrink-0" />
            {!sidebarCollapsed && <span>{label}</span>}
          </NavLink>
        ))}
      </nav>

      {/* Collapse toggle */}
      <div className="p-2 border-t border-slate-700/50">
        <button
          onClick={toggleSidebar}
          className="w-full flex items-center justify-center p-2 rounded-lg text-slate-500 hover:text-slate-300 hover:bg-slate-700/50 transition-colors"
          title={sidebarCollapsed ? 'Expandir sidebar' : 'Recolher sidebar'}
        >
          {sidebarCollapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
        </button>
      </div>
    </aside>
  )
}
