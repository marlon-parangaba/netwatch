import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

/**
 * Garante que o usuário está autenticado.
 * Redireciona para /login caso contrário.
 */
export function useRequireAuth() {
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login', { replace: true })
    }
  }, [isAuthenticated, navigate])

  return isAuthenticated
}

/**
 * Redireciona para /dashboard se já autenticado (usar na página de Login).
 */
export function useRedirectIfAuth() {
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard', { replace: true })
    }
  }, [isAuthenticated, navigate])
}
