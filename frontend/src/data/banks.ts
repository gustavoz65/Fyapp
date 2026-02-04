import type { Bank } from '@/types'

export const banks: Bank[] = [
  {
    id: '1',
    name: 'Banco do Brasil',
    accountType: 'Conta Corrente',
    accountNumber: '12345-6',
    balance: 15420.50,
    connected: true,
    lastSync: '2024-01-27T10:30:00',
    logo: '🏦',
  },
  {
    id: '2',
    name: 'Nubank',
    accountType: 'Conta Digital',
    accountNumber: '98765-4',
    balance: 8750.25,
    connected: true,
    lastSync: '2024-01-27T09:15:00',
    logo: '💜',
  },
  {
    id: '3',
    name: 'Inter',
    accountType: 'Conta Corrente',
    accountNumber: '54321-0',
    balance: 3200.00,
    connected: true,
    lastSync: '2024-01-27T08:00:00',
    logo: '🧡',
  },
]
