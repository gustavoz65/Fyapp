import { Plus, Download, Filter } from 'lucide-react'
import Card from '@/components/Card'
import Table from '@/components/Table'
import Button from '@/components/Button'
import Badge from '@/components/Badge'
import { formatCurrency, formatDate } from '@/lib/utils'
import { transactions } from '@/data/transactions'
import { categories } from '@/data/categories'
import type { Transaction } from '@/types'

export default function Transactions() {
  const columns = [
    {
      key: 'date',
      label: 'Data',
      render: (value: unknown) => formatDate(value as string),
    },
    {
      key: 'description',
      label: 'Descrição',
    },
    {
      key: 'categoryId',
      label: 'Categoria',
      render: (value: unknown) => {
        const category = categories.find(c => c.id === value)
        return category?.name || '-'
      },
    },
    {
      key: 'amount',
      label: 'Valor',
      render: (value: unknown, row: Transaction) => (
        <span className={row.type === 'income' ? 'text-green-600' : 'text-red-600'}>
          {row.type === 'income' ? '+' : '-'} {formatCurrency(value as number)}
        </span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (value: unknown) => {
        const status = value as Transaction['status']
        return (
          <Badge variant={
            status === 'completed' ? 'success' :
            status === 'pending' ? 'warning' : 'error'
          }>
            {status === 'completed' ? 'Concluída' :
             status === 'pending' ? 'Pendente' : 'Cancelada'}
          </Badge>
        )
      },
    },
  ]

  return (
    <div className="space-y-6">
      <Card
        title="Transações"
        subtitle={`${transactions.length} transações registradas`}
        headerAction={
          <div className="flex gap-2">
            <Button variant="secondary" size="sm">
              <Filter className="w-4 h-4 mr-2" />
              Filtrar
            </Button>
            <Button variant="secondary" size="sm">
              <Download className="w-4 h-4 mr-2" />
              Exportar
            </Button>
            <Button size="sm">
              <Plus className="w-4 h-4 mr-2" />
              Nova Transação
            </Button>
          </div>
        }
      >
        <Table<Transaction>
          columns={columns}
          data={transactions}
          searchable
          sortable
        />
      </Card>
    </div>
  )
}
