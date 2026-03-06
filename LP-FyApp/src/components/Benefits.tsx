import { Users, DollarSign, Zap, Headphones } from "lucide-react";
import AnimatedCounter from "./AnimatedCounter";

const benefitsList = [
  {
    number: "1K+",
    title: "Usuários Ativos",
    description: "Pessoas confiando em FyApp para gerenciar suas finanças",
    icon: Users,
    gradient: "from-emerald-500 to-teal-500",
  },
  {
    number: "R$ 10M+",
    title: "Gerenciados",
    description: "Volume total de transações processadas com segurança",
    icon: DollarSign,
    gradient: "from-violet-500 to-purple-500",
  },
  {
    number: "99.9%",
    title: "Disponibilidade",
    description: "Plataforma sempre disponível quando você precisa",
    icon: Zap,
    gradient: "from-blue-500 to-cyan-500",
  },
  {
    number: "24/7",
    title: "Suporte",
    description: "Equipe dedicada pronta para ajudar você",
    icon: Headphones,
    gradient: "from-orange-500 to-amber-500",
  },
];

export default function Benefits() {
  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background with gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-background via-accent/5 to-secondary/5" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        <div className="text-center mb-16 space-y-4 animate-fade-in-up">
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            Por que escolher{" "}
            <span className="text-gradient-accent">FyApp</span>?
          </h2>
          <p className="text-lg text-text-secondary">
            Números que demonstram nossa confiabilidade
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 lg:gap-8">
          {benefitsList.map((benefit, index) => {
            const Icon = benefit.icon;
            return (
              <div
                key={index}
                className="group relative animate-fade-in-up"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                {/* Card */}
                <div className="relative h-full bg-surface rounded-2xl border border-border p-8 text-center transition-all duration-300 hover:shadow-xl hover:-translate-y-2 hover:border-accent/50">
                  {/* Icon */}
                  <div className={`inline-flex p-4 rounded-xl bg-gradient-to-br ${benefit.gradient} shadow-lg mb-4`}>
                    <Icon className="w-6 h-6 text-white" />
                  </div>

                  {/* Number */}
                  <div className="text-5xl font-display font-bold text-primary mb-2 group-hover:scale-110 transition-transform duration-300">
                    {benefit.number}
                  </div>

                  {/* Title */}
                  <h3 className="text-xl font-display font-bold text-primary mb-2">
                    {benefit.title}
                  </h3>

                  {/* Description */}
                  <p className="text-sm text-text-secondary leading-relaxed">
                    {benefit.description}
                  </p>

                  {/* Hover gradient border */}
                  <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-accent/0 via-accent/20 to-accent/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
