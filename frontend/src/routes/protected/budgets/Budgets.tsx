import { Plus } from 'lucide-react'
import Card from '@/components/Card'
import Button from '@/components/Button'
import ProgressBar from '@/components/ProgressBar'
import { formatCurrency } from '@/lib/utils'
import { budgets } from '@/data/budgets'
import { categories } from '@/data/categories'

export default function Budgets() {
  return (
    <div className="space-y-6">
      <Card
        title="Orçamentos"
        subtitle={`${budgets.length} orçamentos ativos`}
        headerAction={
          <Button size="sm">
            <Plus className="w-4 h-4 mr-2" />
            Novo Orçamento
          </Button>
        }
      >
        <div className="space-y-6">
          {budgets.map((budget) => {
            const category = categories.find(c => c.id === budget.categoryId)
            const percentage = (budget.spent / budget.amount) * 100
            const color = percentage >= 100 ? 'danger' : percentage >= 80 ? 'warning' : 'success'

            return (
              <div key={budget.id} className="p-4 border border-gray-200 rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div>
                    <h3 className="font-semibold text-gray-900">{budget.name}</h3>
                    <p className="text-sm text-gray-500">{category?.name}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-600">
                      {formatCurrency(budget.spent)} de {formatCurrency(budget.amount)}
                    </p>
                    <p className="text-xs text-gray-500">{budget.period === 'monthly' ? 'Mensal' : 'Anual'}</p>
                  </div>
                </div>
                <ProgressBar
                  value={budget.spent}
                  max={budget.amount}
                  color={color}
                  showPercent
                />
              </div>
            )
          })}
        </div>
      </Card>
    </div>
  )
}
