import { generateId } from '@/lib/utils'
import type { Transaction } from '@/types'

const startDate = new Date('2023-02-01')
const today = new Date()

const getRandomDate = (): Date => {
  const start = startDate.getTime()
  const end = today.getTime()
  return new Date(start + Math.random() * (end - start))
}

const expenseDescriptions: Record<string, string[]> = {
  '1': ['Supermercado Extra', 'Pão de Açúcar', 'Carrefour', 'Restaurante Japonês', 'iFood', 'Padaria'],
  '2': ['Posto Shell', 'Uber', '99 Pop', 'Metrô', 'Pedágio'],
  '3': ['Aluguel Apartamento', 'Condomínio', 'IPTU 2024', 'Conserto Encanamento'],
  '4': ['Plano Unimed', 'Farmácia', 'Consulta Médica', 'Exames Laboratoriais'],
  '5': ['Mensalidade Faculdade', 'Curso Udemy', 'Livros Amazon'],
  '6': ['Netflix', 'Spotify', 'Cinema', 'Viagem Rio de Janeiro'],
  '7': ['Zara', 'Nike Store', 'C&A'],
  '8': ['Conta de Luz', 'Conta de Água', 'Internet Vivo', 'Celular'],
}

const incomeDescriptions: Record<string, string[]> = {
  '9': ['Salário Janeiro', 'Salário Fevereiro', 'Freelance Website', 'Bônus Anual'],
  '10': ['Dividendos ITSA4', 'Tesouro Selic', 'Rendimento Poupança'],
  '11': ['Venda Notebook', 'Consultoria TI', 'Serviço Design'],
}

const generateTransactions = (): Transaction[] => {
  const transactions: Transaction[] = []
  
  for (let i = 0; i < 500; i++) {
    const isIncome = Math.random() > 0.7
    const categoryIds = isIncome ? ['9', '10', '11'] : ['1', '2', '3', '4', '5', '6', '7', '8']
    const categoryId = categoryIds[Math.floor(Math.random() * categoryIds.length)]
    const descriptions = isIncome ? incomeDescriptions[categoryId] : expenseDescriptions[categoryId]
    
    const minAmount = isIncome ? 500 : 20
    const maxAmount = isIncome ? 8000 : 800
    const amount = Math.random() * (maxAmount - minAmount) + minAmount
    
    const statuses: Transaction['status'][] = ['completed', 'completed', 'completed', 'pending', 'cancelled']
    const status = statuses[Math.floor(Math.random() * statuses.length)]
    
    transactions.push({
      id: generateId(),
      description: descriptions[Math.floor(Math.random() * descriptions.length)],
      amount: parseFloat(amount.toFixed(2)),
      type: isIncome ? 'income' : 'expense',
      categoryId,
      date: getRandomDate().toISOString(),
      status,
      bankId: ['1', '2', '3'][Math.floor(Math.random() * 3)],
      notes: Math.random() > 0.8 ? 'Nota adicional sobre esta transação' : '',
      attachments: [],
      recurring: Math.random() > 0.9,
    })
  }
  
  return transactions.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
}

export const transactions = generateTransactions()
