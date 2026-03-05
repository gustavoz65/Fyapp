import React from 'react';

const featuresList = [
  {
    title: 'Análise Inteligente',
    description: 'Visualize seus gastos em tempo real com gráficos intuitivos, relatórios detalhados e insights que ajudam você a tomar melhores decisões financeiras.',
  },
  {
    title: 'Metas Financeiras',
    description: 'Defina e acompanhe suas metas de economia com precisão. Receba motivação e feedback em tempo real sobre seu progresso.',
  },
  {
    title: 'Múltiplas Contas',
    description: 'Gerencie todas as suas contas bancárias, cartões de crédito e investimentos em um único dashboard centralizado.',
  },
  {
    title: 'Notificações Inteligentes',
    description: 'Receba alertas contextualizados sobre transações importantes, limites de orçamento e oportunidades de economia.',
  },
  {
    title: 'Segurança Máxima',
    description: 'Seus dados são protegidos com criptografia de nível bancário, autenticação de dois fatores e conformidade com regulamentações internacionais.',
  },
  {
    title: 'Sincronização Cloud',
    description: 'Acesse suas finanças de qualquer lugar, em qualquer dispositivo. Sincronização automática e backup seguro de todos os seus dados.',
  },
];

export default function Features() {
  return (
    <section id="features" className="py-32 px-8 max-w-7xl mx-auto">
      {/* Section Header */}
      <div className="text-center mb-16">
        <h2 className="text-5xl font-bold text-[#2C2C2C] mb-4">Recursos Poderosos</h2>
        <p className="text-lg text-[#8B8B8B] max-w-2xl mx-auto">
          Tudo que você precisa para controlar suas finanças em um único lugar
        </p>
      </div>

      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {featuresList.map((feature, index) => (
          <div
            key={index}
            className="bg-white p-8 rounded-2xl border border-[#E5DDD0] transition-all duration-300 hover:shadow-xl hover:-translate-y-2 group relative overflow-hidden"
          >
            {/* Top border gradient on hover */}
            <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-[#7E8C54] to-[#D4A574] transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300 origin-left"></div>

            <h3 className="text-2xl font-bold text-[#2C2C2C] mb-3">{feature.title}</h3>
            <p className="text-[#8B8B8B] leading-relaxed">{feature.description}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
