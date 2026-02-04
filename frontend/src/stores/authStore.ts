import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { AuthState, User } from '@/types'

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: (userData: User) => set({ user: userData, isAuthenticated: true }),
      logout: () => set({ user: null, isAuthenticated: false }),
      updateUser: (userData: Partial<User>) =>
        set((state) => ({ user: state.user ? { ...state.user, ...userData } : null })),
    }),
    {
      name: 'auth-storage',
    }
  )
)
