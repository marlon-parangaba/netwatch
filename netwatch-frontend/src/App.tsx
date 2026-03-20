import { Navigate, Route, Routes } from 'react-router-dom'
import { MainLayout } from '@/components/layout/MainLayout'
import { useAuthStore } from '@/stores/authStore'

// Pages (lazy para code splitting)
import { lazy, Suspense } from 'react'

const Login        = lazy(() => import('@/pages/Login'))
const Dashboard    = lazy(() => import('@/pages/Dashboard'))
const Devices      = lazy(() => import('@/pages/Devices'))
const DeviceDetail = lazy(() => import('@/pages/DeviceDetail'))
const Discovery    = lazy(() => import('@/pages/Discovery'))
const Alerts       = lazy(() => import('@/pages/Alerts'))
const Maps         = lazy(() => import('@/pages/Maps'))
const Settings     = lazy(() => import('@/pages/Settings'))

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function PageLoader() {
  return (
    <div className="flex items-center justify-center h-full min-h-64">
      <div className="w-6 h-6 border-2 border-brand-500 border-t-transparent rounded-full animate-spin" />
    </div>
  )
}

export default function App() {
  return (
    <Suspense fallback={<PageLoader />}>
      <Routes>
        {/* Pública */}
        <Route path="/login" element={<Login />} />

        {/* Protegidas */}
        <Route
          element={
            <ProtectedRoute>
              <MainLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard"    element={<Dashboard />} />
          <Route path="devices"      element={<Devices />} />
          <Route path="devices/:id"  element={<DeviceDetail />} />
          <Route path="discovery"    element={<Discovery />} />
          <Route path="alerts"       element={<Alerts />} />
          <Route path="maps"         element={<Maps />} />
          <Route path="settings"     element={<Settings />} />
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Suspense>
  )
}
