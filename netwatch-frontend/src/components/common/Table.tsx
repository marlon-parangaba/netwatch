import { clsx } from 'clsx'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import type { PaginationMeta } from '@/types'

interface Column<T> {
  key:       string
  header:    string
  render:    (row: T) => React.ReactNode
  className?: string
}

interface TableProps<T> {
  columns:   Column<T>[]
  data:      T[]
  keyFn:     (row: T) => string
  loading?:  boolean
  empty?:    string
}

export function Table<T>({ columns, data, keyFn, loading, empty = 'Nenhum registro encontrado' }: TableProps<T>) {
  return (
    <div className="overflow-x-auto rounded-lg border border-slate-700">
      <table className="w-full text-sm">
        <thead className="bg-slate-800/80">
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className={clsx('px-4 py-3 text-left text-xs font-semibold text-slate-400 uppercase tracking-wider', col.className)}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700/50">
          {loading ? (
            <tr>
              <td colSpan={columns.length} className="px-4 py-12 text-center">
                <div className="flex items-center justify-center gap-2 text-slate-400">
                  <div className="w-5 h-5 border-2 border-brand-500 border-t-transparent rounded-full animate-spin" />
                  <span>Carregando...</span>
                </div>
              </td>
            </tr>
          ) : data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} className="px-4 py-12 text-center text-slate-500">
                {empty}
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr key={keyFn(row)} className="hover:bg-slate-700/30 transition-colors">
                {columns.map((col) => (
                  <td key={col.key} className={clsx('px-4 py-3 text-slate-300', col.className)}>
                    {col.render(row)}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}

interface PaginationProps {
  meta:     PaginationMeta
  onChange: (page: number) => void
}

export function Pagination({ meta, onChange }: PaginationProps) {
  const { page, total_pages, total, limit } = meta
  const from = (page - 1) * limit + 1
  const to   = Math.min(page * limit, total)

  return (
    <div className="flex items-center justify-between px-1 py-2 text-sm text-slate-400">
      <span>{from}–{to} de {total}</span>
      <div className="flex gap-1">
        <button
          onClick={() => onChange(page - 1)}
          disabled={page <= 1}
          className="p-1.5 rounded hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed"
        >
          <ChevronLeft className="w-4 h-4" />
        </button>
        <span className="px-3 py-1">{page} / {total_pages}</span>
        <button
          onClick={() => onChange(page + 1)}
          disabled={page >= total_pages}
          className="p-1.5 rounded hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed"
        >
          <ChevronRight className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}
