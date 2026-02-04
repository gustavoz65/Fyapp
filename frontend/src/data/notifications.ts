import type { Notification } from '@/types'

export const notifications: Notification[] = [
  {
    id: '1',
    type: 'warning',
    title: 'Orçamento quase excedido',
    message: 'Você gastou 96% do seu orçamento de Lazer este mês.',
    date: '2024-01-27T10:00:00',
    read: false,
  },
  {
    id: '2',
    type: 'info',
    title: 'Nova transação sincronizada',
    message: 'Transação de R$ 250,00 detectada no Nubank.',
    date: '2024-01-27T09:30:00',
    read: false,
  },
  {
    id: '3',
    type: 'success',
    title: 'Orçamento cumprido',
    message: 'Parabéns! Você ficou dentro do orçamento de Transporte.',
    date: '2024-01-26T15:00:00',
    read: true,
  },
  {
    id: '4',
    type: 'error',
    title: 'Saldo baixo',
    message: 'Sua conta Inter está com saldo abaixo de R$ 500,00.',
    date: '2024-01-26T10:00:00',
    read: true,
  },
  {
    id: '5',
    type: 'info',
    title: 'Vencimento próximo',
    message: 'Conta de luz vence em 3 dias.',
    date: '2024-01-25T08:00:00',
    read: true,
  },
]
