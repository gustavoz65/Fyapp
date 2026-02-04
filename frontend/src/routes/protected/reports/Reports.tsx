import { useState } from 'react'
import { Download } from 'lucide-react'
import Card from '@/components/Card'
import Button from '@/components/Button'
import Select from '@/components/Select'

export default function Reports() {
  const [reportType, setReportType] = useState('cashflow')

  return (
    <div className="space-y-6">
      <Card
        title="Relatórios"
        subtitle="Gere relatórios detalhados das suas finanças"
        headerAction={
          <Button size="sm">
            <Download className="w-4 h-4 mr-2" />
            Exportar PDF
          </Button>
        }
      >
        <div className="mb-6">
          <Select
            value={reportType}
            onChange={(e) => setReportType(e.target.value)}
            options={[
              { value: 'cashflow', label: 'Fluxo de Caixa' },
              { value: 'dre', label: 'DRE - Demonstração do Resultado' },
              { value: 'balance', label: 'Balanço Patrimonial' },
            ]}
            label="Tipo de Relatório"
          />
        </div>

        <div className="text-center py-20 text-gray-500">
          <p>Selecione um tipo de relatório e configure os filtros para gerar</p>
        </div>
      </Card>
    </div>
  )
}
