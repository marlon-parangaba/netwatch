import { clsx } from 'clsx'
import type { DeviceStatus, AlertSeverity } from '@/types'

type BadgeVariant = 'green' | 'red' | 'yellow' | 'blue' | 'gray' | 'purple'

interface BadgeProps {
  children: React.ReactNode
  variant?: BadgeVariant
  dot?:     boolean
  className?: string
}

const variantClasses: Record<BadgeVariant, string> = {
  green:  'bg-green-400/10  text-green-400  border-green-400/20',
  red:    'bg-red-400/10    text-red-400    border-red-400/20',
  yellow: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/20',
  blue:   'bg-blue-400/10   text-blue-400   border-blue-400/20',
  gray:   'bg-slate-400/10  text-slate-400  border-slate-400/20',
  purple: 'bg-purple-400/10 text-purple-400 border-purple-400/20',
}

export function Badge({ children, variant = 'gray', dot, className }: BadgeProps) {
  return (
    <span className={clsx(
      'inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border',
      variantClasses[variant],
      className
    )}>
      {dot && <span className={clsx('w-1.5 h-1.5 rounded-full', {
        'bg-green-400':  variant === 'green',
        'bg-red-400':    variant === 'red',
        'bg-yellow-400': variant === 'yellow',
        'bg-blue-400':   variant === 'blue',
        'bg-slate-400':  variant === 'gray',
      })} />}
      {children}
    </span>
  )
}

export function StatusBadge({ status }: { status: DeviceStatus }) {
  const map: Record<DeviceStatus, { label: string; variant: BadgeVariant }> = {
    online:  { label: 'Online',  variant: 'green' },
    offline: { label: 'Offline', variant: 'red' },
    warning: { label: 'Warning', variant: 'yellow' },
    unknown: { label: 'Unknown', variant: 'gray' },
  }
  const { label, variant } = map[status]
  return <Badge variant={variant} dot>{label}</Badge>
}

export function SeverityBadge({ severity }: { severity: AlertSeverity }) {
  const map: Record<AlertSeverity, { label: string; variant: BadgeVariant }> = {
    critical: { label: 'Crítico', variant: 'red' },
    high:     { label: 'Alto',    variant: 'red' },
    medium:   { label: 'Médio',   variant: 'yellow' },
    low:      { label: 'Baixo',   variant: 'blue' },
    info:     { label: 'Info',    variant: 'gray' },
  }
  const { label, variant } = map[severity]
  return <Badge variant={variant}>{label}</Badge>
}
