import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { SettingsState } from '@/types'

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      currency: 'BRL',
      language: 'pt-BR',
      notifications: {
        email: true,
        push: true,
        sms: false,
      },
      theme: 'system',
      setCurrency: (currency) => set({ currency }),
      setLanguage: (language) => set({ language }),
      setNotifications: (notifications) => set({ notifications }),
      setTheme: (theme) => set({ theme }),
    }),
    {
      name: 'settings-storage',
    }
  )
)
