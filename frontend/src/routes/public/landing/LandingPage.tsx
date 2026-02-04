import { Link } from 'react-router-dom'
import { 
  BarChart3, 
  TrendingUp, 
  Shield, 
  Smartphone, 
  Users, 
  Zap,
  Check,
  Star,
  ArrowRight
} from 'lucide-react'
import Button from '@/components/Button'

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0a0a0a]">
      <nav className="border-b border-gray-200 dark:border-gray-800">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <BarChart3 className="w-8 h-8 text-primary-600" />
              <span className="text-2xl font-bold text-gray-900 dark:text-white">Cashing</span>
            </div>
            <div className="flex items-center gap-4">
              <Link to="/auth/login">
                <Button variant="ghost">Entrar</Button>
              </Link>
              <Link to="/auth/register">
                <Button>Criar Conta</Button>
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <section className="py-20 bg-gradient-to-b from-primary-50 to-white dark:from-gray-900 dark:to-[#0a0a0a]">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
              Gestão Financeira <span className="text-primary-600">Inteligente</span>
            </h1>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
              Controle suas finanças pessoais e empresariais em um só lugar. 
              Relatórios automáticos, integração bancária e insights em tempo real.
            </p>
            <div className="flex gap-4 justify-center">
              <Link to="/auth/register">
                <Button size="lg">
                  Começar Gratuitamente <ArrowRight className="ml-2 w-5 h-5" />
                </Button>
              </Link>
              <Button variant="secondary" size="lg">
                Ver Demo
              </Button>
            </div>
          </div>
        </div>
      </section>

      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-12">
            Tudo que você precisa para controlar suas finanças
          </h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: BarChart3,
                title: 'Dashboard Intuitivo',
                description: 'Visualize todas as suas finanças em um só lugar com gráficos e métricas em tempo real.',
              },
              {
                icon: TrendingUp,
                title: 'Relatórios Avançados',
                description: 'Gere relatórios de DRE, fluxo de caixa e balanço patrimonial automaticamente.',
              },
              {
                icon: Shield,
                title: 'Segurança Total',
                description: 'Seus dados protegidos com criptografia de ponta a ponta e autenticação em dois fatores.',
              },
              {
                icon: Smartphone,
                title: 'Acesso Mobile',
                description: 'Gerencie suas finanças de qualquer lugar, a qualquer momento.',
              },
              {
                icon: Users,
                title: 'Multi-usuário',
                description: 'Colabore com sua equipe ou família com permissões personalizadas.',
              },
              {
                icon: Zap,
                title: 'Automação',
                description: 'Categorização automática, alertas inteligentes e sincronização bancária.',
              },
            ].map((feature, idx) => (
              <div key={idx} className="p-6 border border-gray-200 dark:border-gray-800 bg-white dark:bg-[#1a1a1a] rounded-xl hover:shadow-lg transition-shadow">
                <feature.icon className="w-12 h-12 text-primary-600 mb-4" />
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{feature.title}</h3>
                <p className="text-gray-600 dark:text-gray-300">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-20 bg-gray-50 dark:bg-[#0f0f0f]">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8 text-center">
            {[
              { number: '10K+', label: 'Usuários Ativos' },
              { number: '99.9%', label: 'Uptime' },
              { number: 'R$ 2B+', label: 'Transacionado' },
              { number: '4.9/5', label: 'Avaliação' },
            ].map((stat, idx) => (
              <div key={idx}>
                <div className="text-4xl font-bold text-primary-600 mb-2">{stat.number}</div>
                <div className="text-gray-600 dark:text-gray-300">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-12">
            O que nossos clientes dizem
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            {[
              {
                name: 'Maria Silva',
                role: 'Empreendedora',
                image: '👩',
                text: 'O Cashing transformou a forma como gerencio as finanças da minha empresa. Relatórios claros e interface intuitiva.',
              },
              {
                name: 'João Santos',
                role: 'Freelancer',
                image: '👨',
                text: 'Finalmente consigo ter uma visão completa das minhas receitas e despesas. Recomendo!',
              },
              {
                name: 'Ana Costa',
                role: 'Contadora',
                image: '👩',
                text: 'Uso com meus clientes. A geração automática de relatórios economiza horas de trabalho.',
              },
            ].map((testimonial, idx) => (
              <div key={idx} className="p-6 bg-white dark:bg-[#1a1a1a] border border-gray-200 dark:border-gray-800 rounded-xl">
                <div className="flex items-center gap-1 mb-4">
                  {[...Array(5)].map((_, i) => (
                    <Star key={i} className="w-5 h-5 fill-yellow-400 text-yellow-400" />
                  ))}
                </div>
                <p className="text-gray-600 dark:text-gray-300 mb-4">{testimonial.text}</p>
                <div className="flex items-center gap-3">
                  <div className="text-4xl">{testimonial.image}</div>
                  <div>
                    <div className="font-semibold text-gray-900 dark:text-white">{testimonial.name}</div>
                    <div className="text-sm text-gray-600 dark:text-gray-400">{testimonial.role}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-20 bg-gray-50 dark:bg-[#0f0f0f]">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-4">
            Escolha o plano ideal para você
          </h2>
          <div className="flex justify-center gap-4 mb-12">
            <Button variant="ghost">Mensal</Button>
            <Button>Anual (20% off)</Button>
          </div>
          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {[
              {
                name: 'Free',
                price: 'R$ 0',
                period: '/mês',
                features: [
                  'Até 100 transações/mês',
                  'Dashboard básico',
                  'Relatórios simples',
                  'Suporte por email',
                ],
              },
              {
                name: 'Pro',
                price: 'R$ 49',
                period: '/mês',
                features: [
                  'Transações ilimitadas',
                  'Dashboard avançado',
                  'Todos os relatórios',
                  'Integração bancária',
                  'Suporte prioritário',
                  'Múltiplos usuários',
                ],
                popular: true,
              },
              {
                name: 'Business',
                price: 'R$ 149',
                period: '/mês',
                features: [
                  'Tudo do Pro',
                  'API de integração',
                  'Customização avançada',
                  'Gerente de conta dedicado',
                  'SLA garantido',
                  'Auditoria e compliance',
                ],
              },
            ].map((plan, idx) => (
              <div
                key={idx}
                className={`p-8 bg-white dark:bg-[#1a1a1a] rounded-xl border-2 ${
                  plan.popular ? 'border-primary-600' : 'border-gray-200 dark:border-gray-800'
                } relative`}
              >
                {plan.popular && (
                  <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
                    <span className="bg-primary-600 text-white px-4 py-1 rounded-full text-sm font-medium">
                      Mais Popular
                    </span>
                  </div>
                )}
                <div className="text-center mb-6">
                  <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">{plan.name}</h3>
                  <div className="flex items-baseline justify-center gap-1">
                    <span className="text-4xl font-bold text-gray-900 dark:text-white">{plan.price}</span>
                    <span className="text-gray-600 dark:text-gray-400">{plan.period}</span>
                  </div>
                </div>
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature, fidx) => (
                    <li key={fidx} className="flex items-start gap-3">
                      <Check className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
                      <span className="text-gray-600 dark:text-gray-300">{feature}</span>
                    </li>
                  ))}
                </ul>
                <Link to="/auth/register">
                  <Button
                    variant={plan.popular ? 'primary' : 'secondary'}
                    className="w-full"
                  >
                    Começar Agora
                  </Button>
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-12">
            Perguntas Frequentes
          </h2>
          <div className="max-w-3xl mx-auto space-y-6">
            {[
              {
                q: 'Como funciona a integração bancária?',
                a: 'Utilizamos tecnologia Open Banking para conectar sua conta de forma segura e automática. Suas credenciais nunca são armazenadas.',
              },
              {
                q: 'Posso usar para pessoa física e jurídica?',
                a: 'Sim! O Cashing foi desenvolvido para atender tanto pessoas físicas quanto jurídicas, com funcionalidades específicas para cada caso.',
              },
              {
                q: 'Meus dados estão seguros?',
                a: 'Absolutamente. Usamos criptografia de ponta a ponta, autenticação em dois fatores e seguimos todas as normas da LGPD.',
              },
              {
                q: 'Posso cancelar a qualquer momento?',
                a: 'Sim, você pode cancelar sua assinatura a qualquer momento, sem multas ou taxas de cancelamento.',
              },
              {
                q: 'Há suporte técnico disponível?',
                a: 'Sim, oferecemos suporte por email para todos os planos, e suporte prioritário para clientes Pro e Business.',
              },
            ].map((faq, idx) => (
              <details key={idx} className="bg-white dark:bg-[#1a1a1a] border border-gray-200 dark:border-gray-800 rounded-lg p-6">
                <summary className="font-semibold text-gray-900 dark:text-white cursor-pointer">{faq.q}</summary>
                <p className="mt-3 text-gray-600 dark:text-gray-300">{faq.a}</p>
              </details>
            ))}
          </div>
        </div>
      </section>

      <footer className="bg-gray-900 dark:bg-[#1a1a1a] text-white py-12 border-t border-gray-800">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <BarChart3 className="w-6 h-6" />
                <span className="text-xl font-bold">Cashing</span>
              </div>
              <p className="text-gray-400 dark:text-gray-300">
                Gestão financeira inteligente para pessoas físicas e jurídicas.
              </p>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Produto</h4>
              <ul className="space-y-2 text-gray-400 dark:text-gray-300">
                <li><a href="#" className="hover:text-white">Funcionalidades</a></li>
                <li><a href="#" className="hover:text-white">Preços</a></li>
                <li><a href="#" className="hover:text-white">Segurança</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Empresa</h4>
              <ul className="space-y-2 text-gray-400 dark:text-gray-300">
                <li><a href="#" className="hover:text-white">Sobre</a></li>
                <li><a href="#" className="hover:text-white">Blog</a></li>
                <li><a href="#" className="hover:text-white">Carreiras</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4">Suporte</h4>
              <ul className="space-y-2 text-gray-400 dark:text-gray-300">
                <li><a href="#" className="hover:text-white">Central de Ajuda</a></li>
                <li><a href="#" className="hover:text-white">Contato</a></li>
                <li><a href="#" className="hover:text-white">Status</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 pt-8 text-center text-gray-400 dark:text-gray-300">
            <p>&copy; 2024 Cashing. Todos os direitos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
