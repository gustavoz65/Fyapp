import { UserPlus, Link2, TrendingUp, ArrowRight } from "lucide-react";

const ONBOARDING_IMG =
  "https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-onboarding-illustration-WgVHrmWdMcQTz6pQ2GqxJK.webp";

const steps = [
  {
    number: 1,
    title: "Crie sua Conta",
    description:
      "Registre-se em segundos com seu email ou telefone. Personalize suas preferências e configure suas notificações.",
    image: ONBOARDING_IMG,
    icon: UserPlus,
    gradient: "from-emerald-500 to-teal-500",
  },
  {
    number: 2,
    title: "Conecte suas Contas",
    description:
      "Integre seus bancos e cartões de forma segura. FyApp sincroniza automaticamente todas as suas transações.",
    image: ONBOARDING_IMG,
    icon: Link2,
    gradient: "from-violet-500 to-purple-500",
  },
  {
    number: 3,
    title: "Comece a Controlar",
    description:
      "Visualize, analise e otimize suas finanças. Defina metas e acompanhe seu progresso em tempo real.",
    image: ONBOARDING_IMG,
    icon: TrendingUp,
    gradient: "from-blue-500 to-cyan-500",
  },
];

export default function HowItWorks() {
  return (
    <section id="how" className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-surface" />
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:14px_24px]" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        {/* Section Header */}
        <div className="text-center mb-16 sm:mb-20 space-y-4 animate-fade-in-up">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-secondary/10 border border-secondary/20">
            <div className="w-2 h-2 rounded-full bg-secondary animate-pulse" />
            <span className="text-sm font-semibold text-secondary">Como Funciona</span>
          </div>
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            Simples e rápido,{" "}
            <span className="text-gradient-multi">em 3 passos</span>
          </h2>
          <p className="text-lg text-text-secondary max-w-2xl mx-auto">
            Começar com FyApp é fácil e intuitivo — você estará no controle das suas finanças em minutos
          </p>
        </div>

        {/* Steps */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12 relative">
          {/* Connection Lines (hidden on mobile) */}
          <div className="hidden md:block absolute top-16 left-1/4 right-1/4 h-0.5 bg-gradient-to-r from-border via-accent to-border" style={{ top: "4rem" }} />

          {steps.map((step, index) => {
            const Icon = step.icon;
            return (
              <div
                key={step.number}
                className="flex flex-col text-center group animate-fade-in-up"
                style={{ animationDelay: `${index * 150}ms` }}
              >
                {/* Step Number Badge */}
                <div className="relative mx-auto mb-8">
                  <div className={`relative w-20 h-20 rounded-2xl bg-gradient-to-br ${step.gradient} shadow-lg flex items-center justify-center group-hover:scale-110 transition-transform duration-300`}>
                    <Icon className="w-10 h-10 text-white" />
                  </div>
                  <div className="absolute -top-2 -right-2 w-8 h-8 rounded-full bg-primary text-white flex items-center justify-center text-sm font-bold shadow-md">
                    {step.number}
                  </div>
                  {/* Glow effect */}
                  <div className={`absolute inset-0 rounded-2xl bg-gradient-to-br ${step.gradient} blur-xl opacity-30 group-hover:opacity-50 transition-opacity duration-300`} />
                </div>

                {/* Content */}
                <h3 className="text-2xl font-display font-bold text-primary mb-4 group-hover:text-accent transition-colors duration-300">
                  {step.title}
                </h3>
                <p className="text-text-secondary leading-relaxed mb-6 flex-1">
                  {step.description}
                </p>

                {/* Image */}
                <div className="relative mt-auto">
                  <div className="absolute -inset-2 bg-gradient-to-r from-accent/20 via-secondary/20 to-accent/20 rounded-xl blur-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                  <div className="relative rounded-xl overflow-hidden border border-border shadow-lg bg-surface group-hover:shadow-xl transition-shadow duration-300">
                    <img
                      src={step.image}
                      alt={`Passo ${step.number} - ${step.title}`}
                      className="w-full h-auto"
                    />
                  </div>
                </div>

                {/* Arrow (except for last step on desktop) */}
                {index < steps.length - 1 && (
                  <div className="hidden md:flex absolute top-16 right-0 translate-x-1/2 items-center justify-center w-8 h-8 rounded-full bg-accent/10 border border-accent/20" style={{ top: "3.5rem" }}>
                    <ArrowRight className="w-4 h-4 text-accent" />
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
