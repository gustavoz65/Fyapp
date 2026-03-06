import { TrendingUp, Target, ArrowRight, Sparkles } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

const highlights = [
  {
    icon: TrendingUp,
    title: "Economia Inteligente",
    description: "Identifique oportunidades de economia automaticamente",
  },
  {
    icon: Target,
    title: "Metas Alcançáveis",
    description: "Defina objetivos e acompanhe seu progresso em tempo real",
  },
];

export default function GrowthSection() {
  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-background via-secondary/5 to-background" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Image - Left Side */}
          <div className="order-2 lg:order-1 relative animate-fade-in-up lg:animate-slide-in-left">
            <div className="relative">
              {/* Glow Effect */}
              <div className="absolute -inset-4 bg-gradient-to-r from-secondary/20 via-accent/20 to-secondary/20 rounded-3xl blur-3xl opacity-50 animate-pulse-slow" />

              {/* Image Container */}
              <div className="relative rounded-2xl overflow-hidden border border-border/50 shadow-2xl bg-background-alt p-8">
                <img
                  src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-growth-illustration-J647hEMMMwFwwztfBg6Kws.webp"
                  alt="Crescimento Financeiro com FyApp"
                  className="w-full rounded-xl"
                />
              </div>

              {/* Floating Stats */}
              <div className="absolute -top-6 -right-6 glass rounded-xl p-4 shadow-xl animate-float hidden lg:block">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-accent/20 flex items-center justify-center">
                    <TrendingUp className="w-5 h-5 text-accent" />
                  </div>
                  <div>
                    <div className="text-xs text-text-tertiary">Crescimento</div>
                    <div className="text-sm font-bold text-primary">+40%</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Content - Right Side */}
          <div className="space-y-8 order-1 lg:order-2 animate-fade-in-up lg:animate-slide-in-right">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent/10 border border-accent/20">
              <Sparkles className="w-4 h-4 text-accent" />
              <span className="text-sm font-semibold text-accent">Crescimento Financeiro</span>
            </div>

            {/* Heading */}
            <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary leading-tight">
              Veja seu dinheiro{" "}
              <span className="text-gradient-accent">crescer</span>
            </h2>

            {/* Description */}
            <p className="text-lg text-text-secondary leading-relaxed">
              Com FyApp, você não apenas controla seus gastos, mas também planeja
              seu futuro financeiro. Nossas ferramentas de análise ajudam você a
              identificar oportunidades de economia e investimento.
            </p>

            {/* Highlight Cards */}
            <div className="space-y-4">
              {highlights.map((highlight, index) => {
                const Icon = highlight.icon;
                return (
                  <div
                    key={index}
                    className="flex items-start gap-4 p-4 rounded-xl bg-surface border border-border hover:border-accent/50 transition-all duration-300 group"
                  >
                    <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-accent/10 flex items-center justify-center group-hover:bg-accent group-hover:scale-110 transition-all duration-300">
                      <Icon className="w-6 h-6 text-accent group-hover:text-white transition-colors duration-300" />
                    </div>
                    <div>
                      <h3 className="font-display font-bold text-primary mb-1">
                        {highlight.title}
                      </h3>
                      <p className="text-sm text-text-secondary">
                        {highlight.description}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Stats */}
            <div className="bg-gradient-to-r from-accent/10 to-secondary/10 rounded-xl p-6 border border-accent/20">
              <p className="text-text-secondary mb-2">
                <span className="text-2xl font-bold text-primary">Centenas de usuários</span>
                <span className="block mt-1">já aumentaram sua saúde financeira em mais de 40% usando FyApp</span>
              </p>
            </div>

            {/* CTA */}
            <a
              href={`${APP_URL}/login`}
              className="group inline-flex items-center gap-2 px-8 py-4 bg-accent text-white rounded-xl font-semibold shadow-lg shadow-accent/30 transition-all duration-300 hover:bg-accent-dark hover:-translate-y-1 hover:shadow-xl hover:shadow-accent/40"
            >
              Começar Sua Jornada
              <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
