import { Plus } from 'lucide-react'
import Card from '@/components/Card'
import Button from '@/components/Button'
import Badge from '@/components/Badge'
import { categories } from '@/data/categories'

export default function Categories() {
  return (
    <div className="space-y-6">
      <Card
        title="Categorias"
        subtitle={`${categories.length} categorias cadastradas`}
        headerAction={
          <Button size="sm">
            <Plus className="w-4 h-4 mr-2" />
            Nova Categoria
          </Button>
        }
      >
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {categories.map((category) => (
            <div
              key={category.id}
              className="p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer"
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-lg flex items-center justify-center text-white text-xl"
                    style={{ backgroundColor: category.color }}
                  >
                    {category.icon === 'utensils' ? '🍴' :
                     category.icon === 'car' ? '🚗' :
                     category.icon === 'home' ? '🏠' :
                     category.icon === 'heart' ? '❤️' :
                     category.icon === 'book' ? '📚' :
                     category.icon === 'gamepad' ? '🎮' :
                     category.icon === 'shirt' ? '👕' :
                     category.icon === 'receipt' ? '🧾' :
                     category.icon === 'briefcase' ? '💼' :
                     category.icon === 'trending-up' ? '📈' :
                     category.icon === 'shopping-bag' ? '🛍️' : '📋'}
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">{category.name}</h3>
                    <p className="text-sm text-gray-500">
                      {category.subcategories?.length || 0} subcategorias
                    </p>
                  </div>
                </div>
                <Badge variant={category.type === 'income' ? 'success' : 'error'}>
                  {category.type === 'income' ? 'Receita' : 'Despesa'}
                </Badge>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  )
}
