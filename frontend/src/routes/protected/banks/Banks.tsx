import { Plus, RefreshCw } from 'lucide-react'
import Card from '@/components/Card'
import Button from '@/components/Button'
import Badge from '@/components/Badge'
import { formatCurrency } from '@/lib/utils'
import { banks } from '@/data/banks'

export default function Banks() {
  return (
    <div className="space-y-6">
      <Card
        title="Contas Bancárias"
        subtitle={`${banks.length} contas conectadas`}
        headerAction={
          <Button size="sm">
            <Plus className="w-4 h-4 mr-2" />
            Conectar Conta
          </Button>
        }
      >
        <div className="space-y-4">
          {banks.map((bank) => (
            <div
              key={bank.id}
              className="p-6 border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="text-4xl">{bank.logo}</div>
                  <div>
                    <h3 className="font-semibold text-gray-900 text-lg">{bank.name}</h3>
                    <p className="text-sm text-gray-600">{bank.accountType} - {bank.accountNumber}</p>
                    <p className="text-xs text-gray-500 mt-1">
                      Última sincronização: {new Date(bank.lastSync).toLocaleString('pt-BR')}
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-bold text-gray-900 mb-2">
                    {formatCurrency(bank.balance)}
                  </p>
                  <div className="flex items-center gap-2">
                    <Badge variant={bank.connected ? 'success' : 'error'}>
                      {bank.connected ? 'Conectado' : 'Desconectado'}
                    </Badge>
                    <Button variant="ghost" size="sm">
                      <RefreshCw className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  )
}
