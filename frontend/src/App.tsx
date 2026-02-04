import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import { useTheme } from './hooks/useTheme'

import LandingPage from './routes/public/landing/LandingPage'
import Login from './routes/public/auth/Login'
import Register from './routes/public/auth/Register'
import ResetPassword from './routes/public/auth/ResetPassword'
import VerifyEmail from './routes/public/auth/VerifyEmail'

import ProtectedLayout from './routes/protected/ProtectedLayout'
import Dashboard from './routes/protected/dashboard/Dashboard'
import Transactions from './routes/protected/transactions/Transactions'
import Reports from './routes/protected/reports/Reports'
import Categories from './routes/protected/categories/Categories'
import Budgets from './routes/protected/budgets/Budgets'
import Banks from './routes/protected/banks/Banks'
import Settings from './routes/protected/settings/Settings'
import Notifications from './routes/protected/notifications/Notifications'

interface ProtectedRouteProps {
  children: React.ReactNode
}

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? <>{children}</> : <Navigate to="/auth/login" replace />
}

function App() {
  useTheme()
  
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/auth/login" element={<Login />} />
        <Route path="/auth/register" element={<Register />} />
        <Route path="/auth/reset-password" element={<ResetPassword />} />
        <Route path="/auth/verify-email" element={<VerifyEmail />} />

        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <ProtectedLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Dashboard />} />
          <Route path="transactions" element={<Transactions />} />
          <Route path="reports" element={<Reports />} />
          <Route path="categories" element={<Categories />} />
          <Route path="budgets" element={<Budgets />} />
          <Route path="banks" element={<Banks />} />
          <Route path="settings" element={<Settings />} />
          <Route path="notifications" element={<Notifications />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
