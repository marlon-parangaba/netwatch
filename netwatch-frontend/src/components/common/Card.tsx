import { HTMLAttributes } from 'react'
import { clsx } from 'clsx'

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

const paddings = { none: '', sm: 'p-3', md: 'p-5', lg: 'p-6' }

export function Card({ padding = 'md', className, children, ...props }: CardProps) {
  return (
    <div
      className={clsx(
        'bg-slate-800 border border-slate-700 rounded-xl shadow-sm',
        paddings[padding],
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}

interface StatCardProps {
  title: string
  value: string | number
  icon:  React.ReactNode
  color?: 'blue' | 'green' | 'red' | 'yellow' | 'gray'
  sub?:  string
}

const colorMap = {
  blue:   'text-blue-400 bg-blue-400/10',
  green:  'text-green-400 bg-green-400/10',
  red:    'text-red-400 bg-red-400/10',
  yellow: 'text-yellow-400 bg-yellow-400/10',
  gray:   'text-slate-400 bg-slate-400/10',
}

export function StatCard({ title, value, icon, color = 'blue', sub }: StatCardProps) {
  return (
    <Card className="flex items-center gap-4">
      <div className={clsx('p-3 rounded-lg', colorMap[color])}>
        {icon}
      </div>
      <div>
        <p className="text-sm text-slate-400">{title}</p>
        <p className="text-2xl font-bold text-slate-100">{value}</p>
        {sub && <p className="text-xs text-slate-500 mt-0.5">{sub}</p>}
      </div>
    </Card>
  )
}
