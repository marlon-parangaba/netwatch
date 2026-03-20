import { useState, FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Activity, Mail, Lock } from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { Input } from '@/components/common/Input'
import { Button } from '@/components/common/Button'

export default function Login() {
  const [email,    setEmail]    = useState('')
  const [password, setPassword] = useState('')
  const { login, isLoading, error, clearError } = useAuthStore()
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    clearError()
    try {
      await login(email, password)
      navigate('/dashboard')
    } catch {
      // erro já setado na store
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="flex flex-col items-center mb-8">
          <div className="w-14 h-14 bg-brand-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-brand-600/20">
            <Activity className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-slate-100">NetWatch</h1>
          <p className="text-slate-400 text-sm mt-1">Monitoramento SNMP para ISPs</p>
        </div>

        {/* Form */}
        <div className="bg-slate-800 border border-slate-700 rounded-2xl p-6 shadow-xl">
          <h2 className="text-lg font-semibold text-slate-100 mb-5">Entrar</h2>

          {error && (
            <div className="mb-4 px-4 py-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="E-mail"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@netwatch.local"
              leftIcon={<Mail className="w-4 h-4" />}
              autoComplete="email"
              required
            />
            <Input
              label="Senha"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              leftIcon={<Lock className="w-4 h-4" />}
              autoComplete="current-password"
              required
            />
            <Button
              type="submit"
              loading={isLoading}
              className="w-full justify-center mt-2"
              size="lg"
            >
              {isLoading ? 'Entrando...' : 'Entrar'}
            </Button>
          </form>
        </div>

        <p className="text-center text-xs text-slate-600 mt-6">
          NetWatch v0.1.0 — Ertel Telecom
        </p>
      </div>
    </div>
  )
}
