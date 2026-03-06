import { BarChart3, Target, Wallet, Bell, Shield, Cloud, TrendingUp, Sparkles } from "lucide-react";

const featuresList = [
  {
    icon: BarChart3,
    title: "Análise Inteligente",
    description:
      "Visualize seus gastos em tempo real com gráficos intuitivos, relatórios detalhados e insights que ajudam você a tomar melhores decisões financeiras.",
    gradient: "from-emerald-500 to-teal-500",
  },
  {
    icon: Target,
    title: "Metas Financeiras",
    description:
      "Defina e acompanhe suas metas de economia com precisão. Receba motivação e feedback em tempo real sobre seu progresso.",
    gradient: "from-violet-500 to-purple-500",
  },
  {
    icon: Wallet,
    title: "Múltiplas Contas",
    description:
      "Gerencie todas as suas contas bancárias, cartões de crédito e investimentos em um único dashboard centralizado.",
    gradient: "from-blue-500 to-cyan-500",
  },
  {
    icon: Bell,
    title: "Notificações Inteligentes",
    description:
      "Receba alertas contextualizados sobre transações importantes, limites de orçamento e oportunidades de economia.",
    gradient: "from-orange-500 to-amber-500",
  },
  {
    icon: Shield,
    title: "Segurança Máxima",
    description:
      "Seus dados são protegidos com criptografia de nível bancário, autenticação de dois fatores e conformidade com regulamentações internacionais.",
    gradient: "from-red-500 to-pink-500",
  },
  {
    icon: Cloud,
    title: "Sincronização Cloud",
    description:
      "Acesse suas finanças de qualquer lugar, em qualquer dispositivo. Sincronização automática e backup seguro de todos os seus dados.",
    gradient: "from-indigo-500 to-blue-500",
  },
];

export default function Features() {
  return (
    <section id="features" className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 bg-gradient-to-b from-background via-surface to-background" />
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:14px_24px]" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        {/* Section Header */}
        <div className="text-center mb-16 sm:mb-20 space-y-4">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent/10 border border-accent/20">
            <Sparkles className="w-4 h-4 text-accent" />
            <span className="text-sm font-semibold text-accent">Recursos</span>
          </div>
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            Tudo que você precisa
            <br />
            <span className="text-gradient-accent">em um só lugar</span>
          </h2>
          <p className="text-lg text-text-secondary max-w-2xl mx-auto">
            Ferramentas poderosas e intuitivas para tornar sua gestão financeira simples e eficiente
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {featuresList.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <div
                key={index}
                className="group relative bg-surface rounded-2xl border border-border p-8 transition-all duration-300 hover:shadow-xl hover:-translate-y-1 hover:border-accent/50"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                {/* Gradient Glow on Hover */}
                <div className="absolute -inset-px bg-gradient-to-r from-accent/0 via-accent/50 to-accent/0 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 blur-sm" />

                {/* Content */}
                <div className="relative space-y-4">
                  {/* Icon */}
                  <div className={`inline-flex p-3 rounded-xl bg-gradient-to-br ${feature.gradient} shadow-lg`}>
                    <Icon className="w-6 h-6 text-white" />
                  </div>

                  {/* Title */}
                  <h3 className="text-xl font-display font-bold text-primary group-hover:text-accent transition-colors duration-300">
                    {feature.title}
                  </h3>

                  {/* Description */}
                  <p className="text-text-secondary leading-relaxed">
                    {feature.description}
                  </p>

                  {/* Hover Arrow */}
                  <div className="flex items-center gap-2 text-accent opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                    <span className="text-sm font-semibold">Saiba mais</span>
                    <TrendingUp className="w-4 h-4" />
                  </div>
                </div>
              </div>
            );
          })}
        </div>

        {/* Bottom CTA */}
        <div className="mt-16 text-center">
          <p className="text-text-secondary mb-4">
            E muito mais recursos sendo desenvolvidos
          </p>
          <div className="flex items-center justify-center gap-2">
            <div className="flex -space-x-2">
              {[1, 2, 3, 4].map((i) => (
                <div
                  key={i}
                  className="w-8 h-8 rounded-full bg-gradient-to-br from-accent to-secondary border-2 border-surface"
                />
              ))}
            </div>
            <span className="text-sm text-text-tertiary">+1.000 usuários ativos</span>
          </div>
        </div>
      </div>
    </section>
  );
}
