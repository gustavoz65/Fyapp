export interface User {
  id: string
  name: string
  email: string
  type: 'individual' | 'business'
  avatar?: string
  plan: 'free' | 'pro' | 'business'
  createdAt: string
}

export interface Category {
  id: string
  name: string
  icon: string
  color: string
  type: 'income' | 'expense'
  subcategories?: Subcategory[]
}

export interface Subcategory {
  id: string
  name: string
}

export interface Transaction extends Record<string, unknown> {
  id: string
  description: string
  amount: number
  type: 'income' | 'expense'
  categoryId: string
  date: string
  status: 'completed' | 'pending' | 'cancelled'
  bankId: string
  notes?: string
  attachments?: string[]
  recurring?: boolean
}

export interface Bank {
  id: string
  name: string
  accountType: string
  accountNumber: string
  balance: number
  connected: boolean
  lastSync: string
  logo: string
}

export interface Budget {
  id: string
  name: string
  categoryId: string
  amount: number
  spent: number
  period: 'monthly' | 'yearly'
  startDate: string
  alertThreshold: number
}

export interface Notification {
  id: string
  type: 'info' | 'success' | 'warning' | 'error'
  title: string
  message: string
  date: string
  read: boolean
}

export interface AuthState {
  user: User | null
  isAuthenticated: boolean
  login: (userData: User) => void
  logout: () => void
  updateUser: (userData: Partial<User>) => void
}

export interface SettingsState {
  currency: 'BRL' | 'USD' | 'EUR'
  language: 'pt-BR' | 'en'
  notifications: {
    email: boolean
    push: boolean
    sms: boolean
  }
  theme: 'light' | 'dark' | 'system'
  setCurrency: (currency: 'BRL' | 'USD' | 'EUR') => void
  setLanguage: (language: 'pt-BR' | 'en') => void
  setNotifications: (notifications: SettingsState['notifications']) => void
  setTheme: (theme: 'light' | 'dark' | 'system') => void
}
