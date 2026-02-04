import { ArrowUpRight, ArrowDownRight, Wallet, TrendingUp } from 'lucide-react'
import Card from '@/components/Card'
import { formatCurrency } from '@/lib/utils'
import { LineChart, Line, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { transactions } from '@/data/transactions'
import { categories } from '@/data/categories'
import { banks } from '@/data/banks'
import Badge from '@/components/Badge'

export default function Dashboard() {
  const totalBalance = banks.reduce((acc, bank) => acc + bank.balance, 0)
  const currentMonth = new Date().getMonth()
  const currentMonthTransactions = transactions.filter(t => {
    const tMonth = new Date(t.date).getMonth()
    return tMonth === currentMonth && t.status === 'completed'
  })

  const income = currentMonthTransactions
    .filter(t => t.type === 'income')
    .reduce((acc, t) => acc + t.amount, 0)
  
  const expenses = currentMonthTransactions
    .filter(t => t.type === 'expense')
    .reduce((acc, t) => acc + t.amount, 0)

  const monthlyData = Array.from({ length: 12 }, (_, i) => {
    const month = i
    const monthTransactions = transactions.filter(t => {
      const tMonth = new Date(t.date).getMonth()
      return tMonth === month && t.status === 'completed'
    })

    return {
      month: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'][i],
      receitas: monthTransactions.filter(t => t.type === 'income').reduce((acc, t) => acc + t.amount, 0),
      despesas: monthTransactions.filter(t => t.type === 'expense').reduce((acc, t) => acc + t.amount, 0),
    }
  })

  const categoryData = categories
    .filter(c => c.type === 'expense')
    .map(cat => ({
      name: cat.name,
      value: currentMonthTransactions
        .filter(t => t.categoryId === cat.id)
        .reduce((acc, t) => acc + t.amount, 0),
      color: cat.color,
    }))
    .filter(c => c.value > 0)
    .slice(0, 6)

  const recentTransactions = transactions.slice(0, 5)

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 mb-1">Saldo Total</p>
              <p className="text-2xl font-bold text-gray-900">{formatCurrency(totalBalance)}</p>
            </div>
            <div className="p-3 bg-primary-100 rounded-lg">
              <Wallet className="w-6 h-6 text-primary-600" />
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 mb-1">Receitas do Mês</p>
              <p className="text-2xl font-bold text-green-600">{formatCurrency(income)}</p>
            </div>
            <div className="p-3 bg-green-100 rounded-lg">
              <ArrowUpRight className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 mb-1">Despesas do Mês</p>
              <p className="text-2xl font-bold text-red-600">{formatCurrency(expenses)}</p>
            </div>
            <div className="p-3 bg-red-100 rounded-lg">
              <ArrowDownRight className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 mb-1">Saldo Líquido</p>
              <p className={`text-2xl font-bold ${income - expenses >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatCurrency(income - expenses)}
              </p>
            </div>
            <div className="p-3 bg-blue-100 rounded-lg">
              <TrendingUp className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </Card>
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        <Card title="Fluxo de Caixa (12 meses)">
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={monthlyData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip formatter={(value: number) => formatCurrency(value)} />
              <Legend />
              <Line type="monotone" dataKey="receitas" stroke="#10b981" strokeWidth={2} name="Receitas" />
              <Line type="monotone" dataKey="despesas" stroke="#ef4444" strokeWidth={2} name="Despesas" />
            </LineChart>
          </ResponsiveContainer>
        </Card>

        <Card title="Despesas por Categoria">
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={categoryData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={(entry) => entry.name}
                outerRadius={100}
                fill="#8884d8"
                dataKey="value"
              >
                {categoryData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip formatter={(value: number) => formatCurrency(value)} />
            </PieChart>
          </ResponsiveContainer>
        </Card>
      </div>

      <Card title="Últimas Transações">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b border-gray-200">
              <tr>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-600">Data</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-600">Descrição</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-600">Categoria</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-600">Valor</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-600">Status</th>
              </tr>
            </thead>
            <tbody>
              {recentTransactions.map((transaction) => {
                const category = categories.find(c => c.id === transaction.categoryId)
                return (
                  <tr key={transaction.id} className="border-b border-gray-100">
                    <td className="py-3 px-4 text-sm text-gray-600">
                      {new Date(transaction.date).toLocaleDateString('pt-BR')}
                    </td>
                    <td className="py-3 px-4 text-sm font-medium text-gray-900">{transaction.description}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{category?.name}</td>
                    <td className={`py-3 px-4 text-sm font-semibold text-right ${
                      transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {transaction.type === 'income' ? '+' : '-'} {formatCurrency(transaction.amount)}
                    </td>
                    <td className="py-3 px-4">
                      <Badge variant={
                        transaction.status === 'completed' ? 'success' :
                        transaction.status === 'pending' ? 'warning' : 'error'
                      }>
                        {transaction.status === 'completed' ? 'Concluída' :
                         transaction.status === 'pending' ? 'Pendente' : 'Cancelada'}
                      </Badge>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
