import { useState } from 'react'
import { User, Bell, Globe, Shield, Palette, CreditCard } from 'lucide-react'
import Card from '@/components/Card'
import Input from '@/components/Input'
import Select from '@/components/Select'
import Button from '@/components/Button'
import { useAuthStore } from '@/stores/authStore'
import { useSettingsStore } from '@/stores/settingsStore'
import { cn } from '@/lib/utils'

const settingsSections = [
  {
    id: 'profile',
    label: 'Perfil',
    icon: User,
  },
  {
    id: 'appearance',
    label: 'Aparência',
    icon: Palette,
  },
  {
    id: 'notifications',
    label: 'Notificações',
    icon: Bell,
  },
  {
    id: 'language',
    label: 'Idioma e Região',
    icon: Globe,
  },
  {
    id: 'security',
    label: 'Segurança',
    icon: Shield,
  },
  {
    id: 'billing',
    label: 'Faturamento',
    icon: CreditCard,
  },
]

export default function Settings() {
  const [activeSection, setActiveSection] = useState('profile')
  const { user } = useAuthStore()
  const { currency, setCurrency, notifications, setNotifications, theme, setTheme } = useSettingsStore()
  const [loading, setLoading] = useState(false)

  const handleSave = () => {
    setLoading(true)
    setTimeout(() => {
      setLoading(false)
    }, 1000)
  }

  return (
    <div className="flex gap-6 h-[calc(100vh-8rem)]">
      {/* Sidebar de navegação */}
      <aside className="w-64 flex-shrink-0 bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4">
        <nav className="space-y-1">
          {settingsSections.map((section) => {
            const Icon = section.icon
            return (
              <button
                key={section.id}
                onClick={() => setActiveSection(section.id)}
                className={cn(
                  'w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left transition-colors',
                  activeSection === section.id
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800'
                )}
              >
                <Icon className="w-5 h-5" />
                <span className="font-medium">{section.label}</span>
              </button>
            )
          })}
        </nav>
      </aside>

      {/* Conteúdo */}
      <div className="flex-1 overflow-y-auto">
        {activeSection === 'profile' && (
          <Card title="Perfil" subtitle="Gerencie suas informações pessoais">
            <div className="space-y-4">
              <Input label="Nome" defaultValue={user?.name} />
              <Input label="Email" type="email" defaultValue={user?.email} disabled />
              <div className="flex justify-end">
                <Button onClick={handleSave} loading={loading}>
                  Salvar Alterações
                </Button>
              </div>
            </div>
          </Card>
        )}

        {activeSection === 'appearance' && (
          <Card title="Aparência" subtitle="Personalize a aparência do sistema">
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                  Tema
                </label>
                <div className="grid grid-cols-3 gap-4">
                  {[
                    { value: 'light', label: 'Claro', icon: '☀️' },
                    { value: 'dark', label: 'Escuro', icon: '🌙' },
                    { value: 'system', label: 'Sistema', icon: '💻' },
                  ].map((option) => (
                    <button
                      key={option.value}
                      onClick={() => setTheme(option.value as 'light' | 'dark' | 'system')}
                      className={cn(
                        'flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all',
                        theme === option.value
                          ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                          : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                      )}
                    >
                      <span className="text-3xl">{option.icon}</span>
                      <span className="font-medium text-gray-900 dark:text-gray-100">
                        {option.label}
                      </span>
                      {theme === option.value && (
                        <span className="text-xs text-primary-600 dark:text-primary-400">
                          Ativo
                        </span>
                      )}
                    </button>
                  ))}
                </div>
                <p className="mt-3 text-sm text-gray-500 dark:text-gray-400">
                  {theme === 'system'
                    ? 'O tema será ajustado automaticamente de acordo com as preferências do seu sistema'
                    : theme === 'dark'
                    ? 'Tema escuro ativado para reduzir o cansaço visual'
                    : 'Tema claro ativado para melhor legibilidade'}
                </p>
              </div>
            </div>
          </Card>
        )}

        {activeSection === 'notifications' && (
          <Card title="Notificações" subtitle="Configure como deseja receber notificações">
            <div className="space-y-4">
              <label className="flex items-center justify-between p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer">
                <div>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    Notificações por Email
                  </span>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    Receba atualizações importantes por email
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.email}
                  onChange={(e) => setNotifications({ ...notifications, email: e.target.checked })}
                  className="w-5 h-5 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-gray-600 dark:bg-gray-700"
                />
              </label>

              <label className="flex items-center justify-between p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer">
                <div>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    Notificações Push
                  </span>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    Receba notificações em tempo real no navegador
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.push}
                  onChange={(e) => setNotifications({ ...notifications, push: e.target.checked })}
                  className="w-5 h-5 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-gray-600 dark:bg-gray-700"
                />
              </label>

              <label className="flex items-center justify-between p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer">
                <div>
                  <span className="font-medium text-gray-900 dark:text-gray-100">
                    Notificações por SMS
                  </span>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    Receba alertas críticos via SMS
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={notifications.sms}
                  onChange={(e) => setNotifications({ ...notifications, sms: e.target.checked })}
                  className="w-5 h-5 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-gray-600 dark:bg-gray-700"
                />
              </label>
            </div>
          </Card>
        )}

        {activeSection === 'language' && (
          <Card title="Idioma e Região" subtitle="Configure preferências de idioma e moeda">
            <div className="space-y-4">
              <Select
                label="Moeda"
                value={currency}
                onChange={(e) => setCurrency(e.target.value as 'BRL' | 'USD' | 'EUR')}
                options={[
                  { value: 'BRL', label: '🇧🇷 Real Brasileiro (R$)' },
                  { value: 'USD', label: '🇺🇸 Dólar Americano ($)' },
                  { value: 'EUR', label: '🇪🇺 Euro (€)' },
                ]}
              />
            </div>
          </Card>
        )}

        {activeSection === 'security' && (
          <Card title="Segurança" subtitle="Proteja sua conta">
            <div className="space-y-6">
              <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
                <h4 className="font-semibold text-yellow-900 dark:text-yellow-200 mb-2">
                  Autenticação de Dois Fatores
                </h4>
                <p className="text-sm text-yellow-700 dark:text-yellow-300 mb-3">
                  Adicione uma camada extra de segurança à sua conta
                </p>
                <Button variant="secondary" size="sm">
                  Configurar 2FA
                </Button>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  Alterar Senha
                </h4>
                <div className="space-y-3">
                  <Input type="password" label="Senha Atual" />
                  <Input type="password" label="Nova Senha" />
                  <Input type="password" label="Confirmar Nova Senha" />
                </div>
                <div className="flex justify-end mt-4">
                  <Button>Alterar Senha</Button>
                </div>
              </div>
            </div>
          </Card>
        )}

        {activeSection === 'billing' && (
          <Card title="Faturamento" subtitle="Gerencie seu plano e pagamentos">
            <div className="space-y-6">
              <div className="p-6 bg-gradient-to-r from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 rounded-xl border border-primary-200 dark:border-primary-800">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      Plano {user?.plan ? user.plan.charAt(0).toUpperCase() + user.plan.slice(1) : 'Free'}
                    </h4>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      Faturado mensalmente
                    </p>
                  </div>
                  <Button variant="secondary">Mudar Plano</Button>
                </div>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">
                  Histórico de Pagamentos
                </h4>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Nenhum pagamento registrado
                </p>
              </div>
            </div>
          </Card>
        )}
      </div>
    </div>
  )
}
